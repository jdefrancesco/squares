#!/usr/bin/env bash

set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

APP_NAME="${APP_NAME:-Squares}"
EXECUTABLE_NAME="${EXECUTABLE_NAME:-squares}"
CMD_PKG="${CMD_PKG:-./cmd/squares}"

BUNDLE_ID="${BUNDLE_ID:-com.example.squares}"
VERSION="${VERSION:-0.1.0}"
BUILD_NUMBER="${BUILD_NUMBER:-${VERSION}}"

DIST_DIR="${DIST_DIR:-${PROJECT_ROOT}/dist}"
APP_DIR="${APP_DIR:-${DIST_DIR}/${APP_NAME}.app}"
CONTENTS_DIR="${APP_DIR}/Contents"
MACOS_DIR="${CONTENTS_DIR}/MacOS"
RESOURCES_DIR="${CONTENTS_DIR}/Resources"

ICON_PNG="${ICON_PNG:-${PROJECT_ROOT}/assets/icon.png}"
ICON_ICNS_NAME="${ICON_ICNS_NAME:-${APP_NAME}.icns}"

GOFLAGS="${GOFLAGS:-}"
LDFLAGS="${LDFLAGS:--s -w}"
CGO_ENABLED_VALUE="${CGO_ENABLED:-1}"

host_arch() {
  uname -m
}

can_lipo() {
  command -v lipo >/dev/null 2>&1
}

build_universal_binary() {
  local out="$1"
  local tmp_dir
  tmp_dir="$(mktemp -d)"

  echo "Building universal binary (arm64 + amd64)…"
  CGO_ENABLED=${CGO_ENABLED_VALUE} GOOS=darwin GOARCH=arm64 go build ${GOFLAGS} -ldflags "${LDFLAGS}" -o "${tmp_dir}/${EXECUTABLE_NAME}-arm64" "${CMD_PKG}"
  CGO_ENABLED=${CGO_ENABLED_VALUE} GOOS=darwin GOARCH=amd64 go build ${GOFLAGS} -ldflags "${LDFLAGS}" -o "${tmp_dir}/${EXECUTABLE_NAME}-amd64" "${CMD_PKG}"

  lipo -create -output "${out}" "${tmp_dir}/${EXECUTABLE_NAME}-arm64" "${tmp_dir}/${EXECUTABLE_NAME}-amd64"
  rm -rf "${tmp_dir}"
}

build_single_arch_binary() {
  local out="$1"
  local arch
  arch="$(host_arch)"

  echo "Building single-arch binary for ${arch}…"
  case "${arch}" in
    arm64) GOOS=darwin GOARCH=arm64 ;;
    x86_64) GOOS=darwin GOARCH=amd64 ;;
    *)
      echo "Unsupported host arch: ${arch}" >&2
      exit 1
      ;;
  esac
  CGO_ENABLED=${CGO_ENABLED_VALUE} go build ${GOFLAGS} -ldflags "${LDFLAGS}" -o "${out}" "${CMD_PKG}"
}

write_info_plist() {
  local plist_path="$1"
  cat >"${plist_path}" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleDevelopmentRegion</key>
  <string>en</string>
  <key>CFBundleExecutable</key>
  <string>${EXECUTABLE_NAME}</string>
  <key>CFBundleIdentifier</key>
  <string>${BUNDLE_ID}</string>
  <key>CFBundleInfoDictionaryVersion</key>
  <string>6.0</string>
  <key>CFBundleName</key>
  <string>${APP_NAME}</string>
  <key>CFBundleDisplayName</key>
  <string>${APP_NAME}</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>CFBundleShortVersionString</key>
  <string>${VERSION}</string>
  <key>CFBundleVersion</key>
  <string>${BUILD_NUMBER}</string>
  <key>LSMinimumSystemVersion</key>
  <string>10.13.0</string>
  <key>NSHighResolutionCapable</key>
  <true/>
  <key>CFBundleIconFile</key>
  <string>${APP_NAME}</string>
</dict>
</plist>
PLIST
}

make_icns_from_png() {
  local png_path="$1"
  local icns_path="$2"

  if [[ ! -f "${png_path}" ]]; then
    echo "No icon found at ${png_path}; skipping .icns generation (the app will use a generic icon)." >&2
    return 0
  fi

  if ! command -v iconutil >/dev/null 2>&1; then
    echo "iconutil not found; skipping .icns generation." >&2
    return 0
  fi

  if ! command -v sips >/dev/null 2>&1; then
    echo "sips not found; skipping .icns generation." >&2
    return 0
  fi

  local iconset_dir
  iconset_dir="$(mktemp -d)"
  iconset_dir="${iconset_dir}/Icon.iconset"
  mkdir -p "${iconset_dir}"

  # Expected input: 1024x1024 PNG with transparency.
  local sizes=(16 32 128 256 512)
  for s in "${sizes[@]}"; do
    sips -z "${s}" "${s}" "${png_path}" --out "${iconset_dir}/icon_${s}x${s}.png" >/dev/null
    local s2=$((s * 2))
    sips -z "${s2}" "${s2}" "${png_path}" --out "${iconset_dir}/icon_${s}x${s}@2x.png" >/dev/null
  done
  sips -z 1024 1024 "${png_path}" --out "${iconset_dir}/icon_512x512@2x.png" >/dev/null

  iconutil -c icns "${iconset_dir}" -o "${icns_path}"
  rm -rf "$(dirname "${iconset_dir}")"
}

main() {
  cd "${PROJECT_ROOT}"
  mkdir -p "${MACOS_DIR}" "${RESOURCES_DIR}"

  echo "Creating app bundle at ${APP_DIR}…"
  rm -rf "${APP_DIR}"
  mkdir -p "${MACOS_DIR}" "${RESOURCES_DIR}"

  write_info_plist "${CONTENTS_DIR}/Info.plist"

  if can_lipo; then
    build_universal_binary "${MACOS_DIR}/${EXECUTABLE_NAME}"
  else
    build_single_arch_binary "${MACOS_DIR}/${EXECUTABLE_NAME}"
  fi

  chmod +x "${MACOS_DIR}/${EXECUTABLE_NAME}"

  make_icns_from_png "${ICON_PNG}" "${RESOURCES_DIR}/${ICON_ICNS_NAME}"

  echo "Done: ${APP_DIR}"
  echo "Tip: open ${APP_NAME}.app from Finder, or run: open \"${APP_DIR}\""
}

main "$@"
