import argparse
import json
from pathlib import Path


def to_table(headers, rows):
    out = []
    out.append("| " + " | ".join(headers) + " |")
    out.append("|" + "|".join(["---"] * len(headers)) + "|")
    for r in rows:
        out.append("| " + " | ".join(str(r.get(h, "")) for h in headers) + " |")
    return "\n".join(out)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--input-dir", required=True)
    parser.add_argument("--output-md", required=True)
    args = parser.parse_args()

    input_dir = Path(args.input_dir)
    rows = json.loads((input_dir / "rows.json").read_text(encoding="utf-8"))
    winners = json.loads((input_dir / "winners.json").read_text(encoding="utf-8"))

    # pick best row per preset+pc class by balanced score
    best = {}
    for r in rows:
        key = (r["preset"], r["pc_class"])
        if key not in best or r["balanced_score"] > best[key]["balanced_score"]:
            best[key] = r

    table_rows = []
    for key in sorted(best.keys(), key=lambda x: (x[1], x[0])):
        r = best[key]
        table_rows.append(
            {
                "preset": r["preset"],
                "pc_class": r["pc_class"],
                "source": r["source"],
                "p95": r["p95_pause_ms"],
                "p99": r["p99_pause_ms"],
                "avg_fps": r["avg_fps"],
                "low1_fps": r["low_1_fps"],
                "full_gc": r["full_gc_count"],
                "startup_ms": r["startup_ms"],
                "score": r["balanced_score"],
            }
        )

    rec_rows = []
    for w in winners:
        rec_rows.append(
            {
                "pc_class": w["pc_class"],
                "recommended_preset": w["winner_preset"],
                "balanced_score": w["balanced_score"],
            }
        )

    # Keep/drop table for JDK 25 strategy.
    flag_rows = [
        {"flag": "-XX:+UseG1GC", "impact": "GC baseline", "risk": "Low", "decision": "Keep"},
        {"flag": "-XX:MaxGCPauseMillis", "impact": "Latency target", "risk": "Medium", "decision": "Keep (50-100)"},
        {"flag": "-XX:InitiatingHeapOccupancyPercent", "impact": "Earlier marking", "risk": "Medium", "decision": "Keep (25-35)"},
        {"flag": "-XX:G1ReservePercent", "impact": "Promotion safety", "risk": "Low", "decision": "Keep (15-20)"},
        {"flag": "-XX:G1HeapRegionSize", "impact": "Region granularity", "risk": "High", "decision": "Drop fixed value"},
        {"flag": "-XX:+AlwaysPreTouch", "impact": "Predictable memory", "risk": "Medium startup", "decision": "Keep on 8+ GB"},
        {"flag": "-XX:+UseStringDeduplication", "impact": "Lower heap pressure", "risk": "Low", "decision": "Keep"},
        {"flag": "Deep G1 refinement flags", "impact": "Niche tuning", "risk": "High portability", "decision": "Drop by default"},
    ]

    md = []
    md.append("# JDK 25 Profile Benchmark Report")
    md.append("")
    md.append("This report is generated from `scripts/perf/*` harness.")
    md.append("")
    md.append("## Preset Metrics (Best Per Source/Class)")
    md.append("")
    md.append(
        to_table(
            ["preset", "pc_class", "source", "p95", "p99", "avg_fps", "low1_fps", "full_gc", "startup_ms", "score"],
            table_rows,
        )
    )
    md.append("")
    md.append("## Recommended Preset Per PC Class")
    md.append("")
    md.append(to_table(["pc_class", "recommended_preset", "balanced_score"], rec_rows))
    md.append("")
    md.append("## JDK 25 Flag Strategy")
    md.append("")
    md.append(to_table(["flag", "impact", "risk", "decision"], flag_rows))
    md.append("")
    md.append("## Scoring Formula")
    md.append("")
    md.append("- `BalancedScore = 0.5 * Latency + 0.4 * Throughput + 0.1 * Stability`")
    md.append("- Latency uses normalized p95/p99 pause.")
    md.append("- Throughput uses normalized Avg FPS and 1% low FPS.")
    md.append("- Stability uses Full GC count and startup time.")

    Path(args.output_md).write_text("\n".join(md) + "\n", encoding="utf-8")


if __name__ == "__main__":
    main()
