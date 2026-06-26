#!/usr/bin/env bash
# Master runner — all automated business-rule verification.
#
# Usage:
#   ./scripts/verify-all-business-rules.sh                         # unit + integration
#   ./scripts/verify-all-business-rules.sh --quick                 # domain/unit only
#   ./scripts/verify-all-business-rules.sh --e2e                   # full + e2e phases + report
#   ./scripts/verify-all-business-rules.sh --e2e --fail-fast       # stop on first failure
#   ./scripts/verify-all-business-rules.sh --e2e --continue-on-error

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FAIL=0
MODE=""
FAIL_FAST=0
CONTINUE_ON_ERROR=1
declare -A PHASE_RESULTS=()

for arg in "$@"; do
  case "$arg" in
    --quick|--e2e) MODE="$arg" ;;
    --fail-fast) FAIL_FAST=1; CONTINUE_ON_ERROR=0 ;;
    --continue-on-error) FAIL_FAST=0; CONTINUE_ON_ERROR=1 ;;
    *) echo "Unknown flag: $arg" >&2; exit 2 ;;
  esac
done

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

on_failure() {
  if [[ "$FAIL_FAST" -eq 1 ]]; then
    if [[ "$MODE" == "--e2e" ]] && [[ "${#PHASE_RESULTS[@]}" -gt 0 ]]; then
      write_e2e_report
    fi
    echo "FAIL-FAST: aborting" >&2
    exit 1
  fi
}

run_script() {
  local name="$1"
  local path="$2"
  echo
  echo "========== $name =========="
  if [[ -x "$path" ]] && "$path"; then
    pass "$name"
    return 0
  fi
  fail "$name"
  on_failure
  return 1
}

run_e2e_phase() {
  local phase="$1"
  local script="$2"
  echo
  echo "========== ${phase} =========="
  if [[ ! -x "$script" ]]; then
    fail "${phase} (missing run.sh)"
    PHASE_RESULTS["$phase"]="SKIP"
    on_failure
    return 1
  fi
  if "$script"; then
    pass "$phase"
    PHASE_RESULTS["$phase"]="PASS"
    return 0
  fi
  fail "$phase"
  PHASE_RESULTS["$phase"]="FAIL"
  on_failure
  return 1
}

write_e2e_report() {
  local overall="PASSED"
  [[ "$FAIL" -ne 0 ]] && overall="FAILED"

  local report_dir="${ROOT}/.local/reports"
  local results_file="${report_dir}/e2e-phase-results.json"
  local report_file="${report_dir}/e2e-last-run.md"
  mkdir -p "$report_dir"

  local flags="--continue-on-error"
  [[ "$FAIL_FAST" -eq 1 ]] && flags="--fail-fast"

  {
    echo '{'
    printf '  "mode": "%s",\n' "${MODE:-full}"
    printf '  "flags": "%s",\n' "$flags"
    echo '  "phases": {'
    local first=1
    for phase in $(printf '%s\n' "${!PHASE_RESULTS[@]}" | sort); do
      [[ "$first" -eq 1 ]] || echo ','
      first=0
      printf '    "%s": "%s"' "$phase" "${PHASE_RESULTS[$phase]}"
    done
    echo
    echo '  },'
    printf '  "overall": "%s"\n' "$overall"
    echo '}'
  } > "$results_file"

  python3 "${ROOT}/scripts/generate-e2e-report.py" \
    "${ROOT}/scripts/business-rules-coverage.json" \
    "$results_file" \
    "$report_file"
  echo "Report: $report_file"
}

echo "=== OpenPresence — business rules verification ==="
echo "Root: $ROOT"
echo "Mode: ${MODE:-full}"
echo "On failure: $( [[ "$FAIL_FAST" -eq 1 ]] && echo fail-fast || echo continue-on-error )"
echo "Matrix: scripts/business-rules-coverage.json"
echo

QUICK_SCRIPTS=(
  "Scaffold|${ROOT}/scripts/verify-scaffold.sh"
  "Geofence BR-020–024|${ROOT}/scripts/verify-geofence.sh"
  "Organization tree|${ROOT}/scripts/verify-organization.sh"
  "Attendance policy|${ROOT}/scripts/verify-attendance-policy.sh"
  "Workforce placement BR-023|${ROOT}/scripts/verify-workforce-placement.sh"
  "Work schedule BR-030–034|${ROOT}/scripts/verify-work-schedule.sh"
  "Punch domain BR-010–015|${ROOT}/scripts/verify-punch.sh"
  "Fraud BR-012–013|${ROOT}/scripts/verify-fraud.sh"
  "Authorization ABAC|${ROOT}/scripts/verify-authorization.sh"
  "Mobile KMP|${ROOT}/scripts/verify-mobile.sh"
  "Biometric stub|${ROOT}/scripts/verify-biometric.sh"
  "Enrollment BR-001–003|${ROOT}/scripts/verify-enrollment.sh"
)

for entry in "${QUICK_SCRIPTS[@]}"; do
  IFS='|' read -r name path <<< "$entry"
  run_script "$name" "$path" || [[ "$FAIL_FAST" -eq 0 ]]
done

if [[ "$MODE" == "--quick" ]]; then
  echo
  if [[ "$FAIL" -eq 0 ]]; then
    echo "=== Quick verification: ALL PASSED ==="
    exit 0
  fi
  echo "=== Quick verification: FAILED ==="
  exit 1
fi

INTEGRATION_SCRIPTS=(
  "RLS multi-tenant|${ROOT}/scripts/verify-rls.sh"
  "SubmitPunch use case|${ROOT}/scripts/verify-punch-usecase.sh"
)

for entry in "${INTEGRATION_SCRIPTS[@]}"; do
  IFS='|' read -r name path <<< "$entry"
  run_script "$name" "$path" || [[ "$FAIL_FAST" -eq 0 ]]
done

if [[ -d "${ROOT}/models" ]] && [[ -f "${ROOT}/models/auraface.onnx" ]]; then
  echo
  echo "========== Biometric ONNX (optional) =========="
  if ONNX_MODELS_PATH="${ROOT}/models" "${ROOT}/scripts/verify-biometric.sh"; then
    pass "Biometric ONNX"
  else
    fail "Biometric ONNX"
    on_failure
  fi
else
  echo
  echo "SKIP: Biometric ONNX — run ./scripts/download-models.sh first"
fi

if [[ "$MODE" == "--e2e" ]]; then
  E2E_DIR="${ROOT}/.local/phases/e2e-testing"
  E2E_PHASES=(
    "e2e-01-enrollment"
    "e2e-02-punch-core"
    "e2e-03-geofence"
    "e2e-04-fraud"
    "e2e-05-work-schedule"
    "e2e-06-authorization"
    "e2e-07-rls-security"
    "e2e-08-biometric-stack"
    "e2e-09-submit-punch-stack"
    "e2e-10-mobile-kmp"
    "e2e-11-full-uc001"
  )
  for phase in "${E2E_PHASES[@]}"; do
    run_e2e_phase "$phase" "${E2E_DIR}/${phase}/run.sh" || [[ "$FAIL_FAST" -eq 0 ]]
  done
  write_e2e_report
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Business rules verification: ALL PASSED ==="
  exit 0
fi
echo "=== Business rules verification: FAILED ==="
exit 1
