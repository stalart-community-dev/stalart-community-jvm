import argparse
import csv
import json
from pathlib import Path
from statistics import mean


PC_CLASSES = {
    "low_end": {"cpu_score": 0.65, "mem_score": 0.70},
    "mid_range": {"cpu_score": 0.85, "mem_score": 0.90},
    "high_end": {"cpu_score": 1.00, "mem_score": 1.00},
}


def load_configs(config_dir: Path):
    configs = {}
    for path in sorted(config_dir.glob("*.json")):
        with path.open("r", encoding="utf-8") as f:
            configs[path.stem] = json.load(f)
    return configs


def synthetic_metrics(cfg: dict, pc: str):
    # Calibrated estimator (relative score model) for balanced comparison.
    pc_base = PC_CLASSES[pc]
    pause = cfg.get("max_gc_pause_millis", 80)
    new_size = cfg.get("g1_new_size_percent", 20)
    ihop = cfg.get("initiating_heap_occupancy_percent", 30)
    p_gc = cfg.get("parallel_gc_threads", 6)
    c_gc = cfg.get("conc_gc_threads", 3)
    pretouch = 1 if cfg.get("pre_touch", False) else 0
    dedup = 1 if cfg.get("use_string_deduplication", False) else 0
    heap = cfg.get("heap_size_gb", 8)

    # Latency model: lower is better.
    p95_pause_ms = (
        0.72 * pause
        + 0.28 * max(25, pause - (ihop - 20) * 0.6)
        + (new_size - 20) * 0.9
        + (10 - p_gc) * 0.8
        + (5 - c_gc) * 0.6
        + (0 if pretouch else 3.0)
    )
    p95_pause_ms /= pc_base["cpu_score"]
    p99_pause_ms = p95_pause_ms * (1.28 + (0.04 if new_size > 30 else 0.0))

    # Throughput / FPS proxy.
    fps = (
        120 * pc_base["cpu_score"]
        + (p_gc * 1.8)
        + (c_gc * 1.2)
        - (100 - pause) * 0.08
        - max(0, 25 - ihop) * 0.45
        + (3 if dedup else 0)
        + (2 if pretouch else 0)
        + (heap - 8) * 0.6
    )
    fps = max(40, fps)
    low1_fps = fps * (1 - min(0.28, p99_pause_ms / 420))

    full_gc = 0 if ihop >= 20 and cfg.get("g1_reserve_percent", 15) >= 15 else 1
    startup_ms = (
        2600
        + (220 if pretouch else -40)
        + (cfg.get("reserved_code_cache_size_mb", 400) - 400) * 0.3
        - p_gc * 16
    ) / pc_base["cpu_score"]
    gc_cpu_pct = min(28, 8 + (100 - pause) * 0.13 + max(0, 25 - ihop) * 0.2)

    return {
        "p95_pause_ms": round(p95_pause_ms, 2),
        "p99_pause_ms": round(p99_pause_ms, 2),
        "avg_fps": round(fps, 2),
        "low_1_fps": round(low1_fps, 2),
        "full_gc_count": int(full_gc),
        "startup_ms": round(startup_ms, 2),
        "gc_cpu_pct": round(gc_cpu_pct, 2),
    }


def normalize_higher(values):
    lo, hi = min(values), max(values)
    if hi == lo:
        return [1.0 for _ in values]
    return [(v - lo) / (hi - lo) for v in values]


def normalize_lower(values):
    lo, hi = min(values), max(values)
    if hi == lo:
        return [1.0 for _ in values]
    return [(hi - v) / (hi - lo) for v in values]


def compute_scores(rows):
    p95 = [r["p95_pause_ms"] for r in rows]
    p99 = [r["p99_pause_ms"] for r in rows]
    fps = [r["avg_fps"] for r in rows]
    low1 = [r["low_1_fps"] for r in rows]
    fullgc = [r["full_gc_count"] for r in rows]
    startup = [r["startup_ms"] for r in rows]

    p95_n = normalize_lower(p95)
    p99_n = normalize_lower(p99)
    fps_n = normalize_higher(fps)
    low1_n = normalize_higher(low1)
    fullgc_n = normalize_lower(fullgc)
    startup_n = normalize_lower(startup)

    for i, row in enumerate(rows):
        latency = 0.65 * p95_n[i] + 0.35 * p99_n[i]
        throughput = 0.65 * fps_n[i] + 0.35 * low1_n[i]
        stability = 0.7 * fullgc_n[i] + 0.3 * startup_n[i]
        score = 100 * (0.5 * latency + 0.4 * throughput + 0.1 * stability)
        row["latency_score"] = round(100 * latency, 2)
        row["throughput_score"] = round(100 * throughput, 2)
        row["stability_score"] = round(100 * stability, 2)
        row["balanced_score"] = round(score, 2)


def load_real_rows(real_csv: Path):
    if not real_csv.exists():
        return []
    rows = []
    with real_csv.open("r", encoding="utf-8", newline="") as f:
        reader = csv.DictReader(f)
        for r in reader:
            rows.append(
                {
                    "preset": r["preset"],
                    "pc_class": r["pc_class"],
                    "source": "real",
                    "p95_pause_ms": float(r["p95_pause_ms"]),
                    "p99_pause_ms": float(r["p99_pause_ms"]),
                    "avg_fps": float(r["avg_fps"]),
                    "low_1_fps": float(r["low_1_fps"]),
                    "full_gc_count": int(r["full_gc_count"]),
                    "startup_ms": float(r["startup_ms"]),
                    "gc_cpu_pct": float(r.get("gc_cpu_pct", 0)),
                }
            )
    return rows


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--mode", choices=["synthetic", "real", "both"], default="both")
    parser.add_argument("--config-dir", required=True)
    parser.add_argument("--real-csv", required=True)
    parser.add_argument("--out-dir", required=True)
    args = parser.parse_args()

    out_dir = Path(args.out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)

    rows = []
    if args.mode in ("synthetic", "both"):
        configs = load_configs(Path(args.config_dir))
        for preset, cfg in configs.items():
            for pc_class in PC_CLASSES:
                m = synthetic_metrics(cfg, pc_class)
                rows.append({"preset": preset, "pc_class": pc_class, "source": "synthetic", **m})

    if args.mode in ("real", "both"):
        rows.extend(load_real_rows(Path(args.real_csv)))

    # score per source+pc_class
    grouped = {}
    for r in rows:
        key = (r["source"], r["pc_class"])
        grouped.setdefault(key, []).append(r)
    for _, group_rows in grouped.items():
        compute_scores(group_rows)

    # aggregate winner per pc class by averaged balanced score
    winners = []
    for pc_class in PC_CLASSES:
        by_preset = {}
        by_preset_fullgc = {}
        for r in rows:
            if r["pc_class"] != pc_class:
                continue
            by_preset.setdefault(r["preset"], []).append(r["balanced_score"])
            by_preset_fullgc.setdefault(r["preset"], []).append(r["full_gc_count"])
        if not by_preset:
            continue
        # Safety gate for default strategy:
        # if any preset can deliver 0 FullGC for this class, consider only those presets.
        zero_fullgc_presets = [
            p for p, vals in by_preset_fullgc.items() if mean(vals) == 0
        ]
        scored_pool = (
            {p: by_preset[p] for p in zero_fullgc_presets}
            if zero_fullgc_presets
            else by_preset
        )
        scored = [(p, mean(vals)) for p, vals in scored_pool.items()]
        scored.sort(key=lambda x: x[1], reverse=True)
        winners.append({"pc_class": pc_class, "winner_preset": scored[0][0], "balanced_score": round(scored[0][1], 2)})

    rows_path = out_dir / "rows.json"
    winners_path = out_dir / "winners.json"
    csv_path = out_dir / "rows.csv"

    with rows_path.open("w", encoding="utf-8") as f:
        json.dump(rows, f, ensure_ascii=False, indent=2)
    with winners_path.open("w", encoding="utf-8") as f:
        json.dump(winners, f, ensure_ascii=False, indent=2)

    headers = [
        "preset", "pc_class", "source", "p95_pause_ms", "p99_pause_ms", "avg_fps",
        "low_1_fps", "full_gc_count", "startup_ms", "gc_cpu_pct",
        "latency_score", "throughput_score", "stability_score", "balanced_score",
    ]
    with csv_path.open("w", encoding="utf-8", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=headers)
        writer.writeheader()
        for r in rows:
            writer.writerow(r)

    template_path = out_dir / "real-runs.template.csv"
    with template_path.open("w", encoding="utf-8", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(
            [
                "preset",
                "pc_class",
                "p95_pause_ms",
                "p99_pause_ms",
                "avg_fps",
                "low_1_fps",
                "full_gc_count",
                "startup_ms",
                "gc_cpu_pct",
            ]
        )


if __name__ == "__main__":
    main()
