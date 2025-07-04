APP_NAME=ScheduledNotifier
VERSION=1.0.0
BUILD_DIR=build
RES_DIR=res
MAC_APP=$(BUILD_DIR)/mac/$(APP_NAME).app
WIN_EXE=$(BUILD_DIR)/windows/$(APP_NAME).exe

.PHONY: clean all mac windows

all: clean mac windows

mac: $(MAC_APP)

windows: $(WIN_EXE)

$(MAC_APP):
	@echo "Building macOS application..."
	@mkdir -p $(MAC_APP)/Contents/MacOS
	@mkdir -p $(MAC_APP)/Contents/Resources
	@GOOS=darwin GOARCH=arm64 go build -o $(MAC_APP)/Contents/MacOS/$(APP_NAME) .
	@chmod +x $(MAC_APP)/Contents/MacOS/$(APP_NAME)
	@cp -r $(RES_DIR)/$(APP_NAME).app/Contents/Info.plist $(MAC_APP)/Contents/
	@echo "macOS app built at $(MAC_APP)"

$(WIN_EXE):
	@echo "Building Windows executable..."
	@mkdir -p $(BUILD_DIR)/windows
	@GOOS=windows GOARCH=amd64 go build -o $(WIN_EXE) .
	@echo "Windows executable built at $(WIN_EXE)"

clean:
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned build directory"

package-mac: mac
	@echo "Packaging macOS app..."
	create-dmg \
		--volname "ScheduledNotifier Installer" \
		--window-pos 200 120 \
		--window-size 600 300 \
		--icon-size 100 \
		--icon "ScheduledNotifier.app" 150 150 \
		--hide-extension "ScheduledNotifier.app" \
		--app-drop-link 450 150 \
		"$(BUILD_DIR)/ScheduledNotifier.dmg" \
		 $(MAC_APP)
	@echo "DMG created at $(BUILD_DIR)/$(APP_NAME)-$(VERSION).dmg"

package-win: windows
	@echo "Packaging Windows executable..."
	@zip -j $(BUILD_DIR)/$(APP_NAME)-$(VERSION)-windows.zip $(WIN_EXE)
	@echo "ZIP created at $(BUILD_DIR)/$(APP_NAME)-$(VERSION)-windows.zip"

package: package-mac package-win