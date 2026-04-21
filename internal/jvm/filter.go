package jvm

import "strings"

var exactRemove = map[string]struct{}{
	"-XX:-PrintCommandLineFlags": {},
	"-XX:+UseG1GC":               {},
	"-XX:-UseG1GC":               {},
	"-XX:+UseZGC":                {},
	"-XX:-UseZGC":                {},
	"-XX:+UseShenandoahGC":       {},
	"-XX:-UseShenandoahGC":       {},
	"-XX:+UseParallelGC":         {},
	"-XX:+UseSerialGC":           {},
	"-XX:+UseCompressedOops":     {},
	"-XX:-UseCompressedOops":     {},
	"-XX:+PerfDisableSharedMem":  {},
	"-XX:-PerfDisableSharedMem":  {},
	"-XX:+UseBiasedLocking":      {},
	"-XX:-UseBiasedLocking":      {},
	"-XX:+UseStringDeduplication": {},
	"-XX:-UseStringDeduplication": {},
	"-XX:+UseNUMA":    {},
	"-XX:-UseNUMA":    {},
	"-XX:+UseDynamicNumberOfGCThreads":   {},
	"-XX:-UseDynamicNumberOfGCThreads":   {},
	"-XX:+AlwaysActAsServerClassMachine": {},
	"-XX:-AlwaysActAsServerClassMachine": {},
	"-XX:+UseXMMForArrayCopy": {},
	"-XX:-UseXMMForArrayCopy": {},
	"-XX:+UseFPUForSpilling": {},
	"-XX:-UseFPUForSpilling": {},
	"-XX:-DontCompileHugeMethods": {},
	"-XX:+DontCompileHugeMethods": {},
	"-XX:+AlwaysPreTouch":          {},
	"-XX:-AlwaysPreTouch":          {},
	"-XX:+ParallelRefProcEnabled":  {},
	"-XX:-ParallelRefProcEnabled":  {},
	"-XX:+DisableExplicitGC":       {},
	"-XX:-DisableExplicitGC":       {},
	"-XX:+G1UseAdaptiveIHOP": {},
	"-XX:-G1UseAdaptiveIHOP": {},
	"-XX:+UnlockExperimentalVMOptions": {},
	"-XX:-UnlockExperimentalVMOptions": {},
	"-XX:+UnlockDiagnosticVMOptions": {},
	"-XX:-UnlockDiagnosticVMOptions": {},
	"-XX:+UseThreadPriorities": {},
	"-XX:-UseThreadPriorities": {},
	"-XX:+UseCounterDecay":   {},
	"-XX:-UseCounterDecay":   {},
	"-XX:+UseLargePages":     {},
	"-XX:-UseLargePages":     {},
}

// java25ExactIncompatible removes legacy flags that are unsupported on
// modern HotSpot and cause fatal startup errors with JDK 25.
var java25ExactIncompatible = map[string]struct{}{
	"-XX:+UseConcMarkSweepGC":   {},
	"-XX:-UseConcMarkSweepGC":   {},
	"-XX:+UseParNewGC":          {},
	"-XX:-UseParNewGC":          {},
	"-XX:+UseBiasedLocking":     {},
	"-XX:-UseBiasedLocking":     {},
	"-XX:+AggressiveOpts":       {},
	"-XX:-AggressiveOpts":       {},
	"-XX:+UseFastAccessorMethods": {},
	"-XX:-UseFastAccessorMethods": {},
	"-XX:+UnlockCommercialFeatures": {},
	"-XX:-UnlockCommercialFeatures": {},
	"-XX:+FlightRecorder":          {},
	"-XX:-FlightRecorder":          {},
}

var prefixRemove = []string{
	"-XX:MaxGCPauseMillis=",
	"-XX:MetaspaceSize=",
	"-XX:MaxMetaspaceSize=",
	"-XX:G1HeapRegionSize=",
	"-XX:G1NewSizePercent=",
	"-XX:G1MaxNewSizePercent=",
	"-XX:G1ReservePercent=",
	"-XX:G1HeapWastePercent=",
	"-XX:G1MixedGCCountTarget=",
	"-XX:InitiatingHeapOccupancyPercent=",
	"-XX:G1MixedGCLiveThresholdPercent=",
	"-XX:G1RSetUpdatingPauseTimePercent=",
	"-XX:G1SATBBufferEnqueueingThresholdPercent=",
	"-XX:G1ConcRSHotCardLimit=",
	"-XX:G1ConcRefinementServiceIntervalMillis=",
	"-XX:GCTimeRatio=",
	"-XX:SurvivorRatio=",
	"-XX:MaxTenuringThreshold=",
	"-XX:ParallelGCThreads=",
	"-XX:ConcGCThreads=",
	"-XX:SoftRefLRUPolicyMSPerMB=",
	"-XX:ReservedCodeCacheSize=",
	"-XX:NonNMethodCodeHeapSize=",
	"-XX:ProfiledCodeHeapSize=",
	"-XX:NonProfiledCodeHeapSize=",
	"-XX:MaxInlineLevel=",
	"-XX:FreqInlineSize=",
	"-XX:InlineSmallCode=",
	"-XX:MaxNodeLimit=",
	"-XX:NodeLimitFudgeFactor=",
	"-XX:NmethodSweepActivity=",
	"-XX:AllocatePrefetchStyle=",
	"-XX:LargePageSizeInBytes=",
	"-XX:AutoBoxCacheMax=",
	"-XX:ThreadPriorityPolicy=",
	"-XX:CompileThresholdScaling=",
	"-Dsun.reflect.inflationThreshold=",
	"-Djdk.reflect.useDirectMethodHandleOnly=",
	"-Dio.netty.jfr.enabled=",
	"-Dio.netty.noUnsafe=",
	"-Dio.netty.tryReflectionSetAccessible=",
	"-Djdk.attach.allowAttachSelf=",
	"--sun-misc-unsafe-memory-access=",
	"-XX:SoftMaxHeapSize=",
	"-Xms",
	"-Xmx",
}

// jvmFlagTakesNextArg reports flags whose value is the following argv
// element (JDK 9+ module system, agents, etc.). Without this, values
// like "java.base/java.lang=ALL-UNNAMED" after "--add-opens" were parsed
// as the main class, breaking the launcher ("--add-opens requires modules").
func jvmFlagTakesNextArg(flag string) bool {
	switch flag {
	case "--add-opens", "--add-exports", "--add-reads", "--patch-module",
		"--upgrade-module-path", "--module-path", "--limit-modules", "--add-modules",
		"--module", "-m", "-p",
		"-javaagent", "-agentlib", "-agentpath":
		return true
	default:
		return false
	}
}

// splitArgs partitions the launcher's argv into JVM flags, the main class,
// and arguments passed to main().
func splitArgs(args []string) (jvm []string, mainClass string, app []string) {
	for i := 0; i < len(args); {
		a := args[i]
		// After "-jar" <file>, the launcher passes app args only — they may
		// start with "-" (e.g. "--version") and must not be parsed as JVM flags.
		if a == "-jar" {
			jvm = append(jvm, a)
			i++
			if i < len(args) {
				jvm = append(jvm, args[i])
				i++
			}
			return jvm, "", args[i:]
		}
		if a == "-classpath" || a == "-cp" {
			jvm = append(jvm, a)
			i++
			if i < len(args) {
				jvm = append(jvm, args[i])
			}
			i++
			continue
		}
		if jvmFlagTakesNextArg(a) {
			jvm = append(jvm, a)
			i++
			if i < len(args) {
				jvm = append(jvm, args[i])
			}
			i++
			continue
		}
		if strings.HasPrefix(a, "-") {
			jvm = append(jvm, a)
			i++
			continue
		}
		mainClass = a
		app = args[i+1:]
		return
	}
	return
}

func shouldRemove(arg string) bool {
	if _, ok := exactRemove[arg]; ok {
		return true
	}
	for _, p := range prefixRemove {
		if strings.HasPrefix(arg, p) {
			return true
		}
	}
	return false
}

func isJava25Incompatible(arg string) bool {
	if _, ok := java25ExactIncompatible[arg]; ok {
		return true
	}
	for _, p := range []string{"-XX:PermSize=", "-XX:MaxPermSize=", "--illegal-access="} {
		if strings.HasPrefix(arg, p) {
			return true
		}
	}
	return false
}

// IsLikelyGameLaunch reports whether argv looks like the actual Minecraft
// client process (not Gravit bootstrap / JavaAgent restart stage).
func IsLikelyGameLaunch(orig []string) bool {
	joined := strings.Join(orig, " ")
	if strings.Contains(joined, "net.minecraft.client.main") ||
		strings.Contains(joined, "net.minecraft.client.Main") ||
		strings.Contains(joined, "net.minecraft.launchwrapper.Launch") ||
		strings.Contains(joined, "cpw.mods") ||
		strings.Contains(joined, "GradleStart") ||
		strings.Contains(joined, "--gameDir") ||
		strings.Contains(joined, "--assetsDir") ||
		strings.Contains(joined, "--version") {
		return true
	}
	return false
}

// FilterArgs strips launcher-injected flags that conflict with ours,
// then splices the generated flags back in, preserving the original
// main class and app arguments.
func FilterArgs(orig, injected []string) []string {
	jvmArgs, mainClass, app := splitArgs(orig)

	filtered := make([]string, 0, len(jvmArgs))
	for _, a := range jvmArgs {
		if !shouldRemove(a) {
			filtered = append(filtered, a)
		}
	}

	result := make([]string, 0, len(filtered)+len(injected)+1+len(app))
	result = append(result, filtered...)
	result = append(result, injected...)
	if mainClass != "" {
		result = append(result, mainClass)
	}
	return append(result, app...)
}

// InjectArgs keeps original JVM args and appends injected args before main class.
// Useful for compatibility fallbacks where stripping launcher arguments is risky.
func InjectArgs(orig, injected []string) []string {
	jvmArgs, mainClass, app := splitArgs(orig)
	result := make([]string, 0, len(jvmArgs)+len(injected)+1+len(app))
	result = append(result, jvmArgs...)
	result = append(result, injected...)
	if mainClass != "" {
		result = append(result, mainClass)
	}
	return append(result, app...)
}

// StripJava25IncompatibleArgs removes known legacy VM options that fail on JDK 25.
// It preserves the launch order and only strips fatal incompatibilities.
func StripJava25IncompatibleArgs(orig []string) []string {
	jvmArgs, mainClass, app := splitArgs(orig)
	filtered := make([]string, 0, len(jvmArgs))
	for i := 0; i < len(jvmArgs); i++ {
		a := jvmArgs[i]
		if isJava25Incompatible(a) {
			continue
		}
		filtered = append(filtered, a)
	}
	result := make([]string, 0, len(filtered)+1+len(app))
	result = append(result, filtered...)
	if mainClass != "" {
		result = append(result, mainClass)
	}
	return append(result, app...)
}
