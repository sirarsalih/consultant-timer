# Consultant Timer

Consultant Timer is a lightweight, local-first Windows desktop application for simple daily consultant time tracking.

## Features

### Time tracking

- **START** begins tracking working time.
- **PAUSE** pauses tracking without losing accumulated time.
- **STOP** stops tracking and clears today's accumulated time.
- Time continues accurately across multiple Start/Pause sessions during the same day.
- Saved time survives closing the program and restarting Windows.
- The timer always starts **paused** after reopening the application.

### Two simultaneous time formats

Both formats are displayed at the same time:

- **Decimal hours:** `7.50 hours`
- **Hours/minutes:** `07:30`
- Both values can be selected and copied.

### Automatic idle detection

- Detects **5 minutes of Windows inactivity**.
- Automatically pauses the timer when the idle threshold is reached.
- Removes the detected idle period instead of counting it as working time.
- Displays that the timer was paused because of inactivity.

### Local and offline operation

- No account required.
- No cloud service.
- No internet connection required.
- Time data is stored locally on the computer.
- Daily data persists between application sessions.
- Previous days remain stored locally when a new day begins.

### Start with Windows

- Includes a **“Start automatically when Windows starts”** checkbox.
- Checkbox and text are centered in the interface.
- Automatic startup can be enabled or disabled directly from Consultant Timer.
- Uses the current Windows user's startup configuration.
- When Windows starts, Consultant Timer launches **paused** rather than automatically counting time.
- The startup preference is remembered.

### Windows interface

- Native, lightweight Windows desktop application.
- Large and simple daily-time display.
- **Green START** button.
- **Yellow PAUSE** button.
- **Red STOP** button.
- Buttons automatically enable or disable depending on the current timer state.
- Hand/pointer cursor when hovering over action buttons.
- Displays the current state, such as **Running**, **Paused**, or **Stopped**.
- Shows that idle detection is set to **5 minutes** and that data is stored locally.

### System tray

- Consultant Timer includes a Windows system-tray icon.
- Right-clicking the tray icon provides:
  - **Start**
  - **Pause**
  - **Stop / Clear today**
  - **Show**
  - **Exit**
- Double-clicking the tray icon shows the main window.
- The tray tooltip indicates the current timer state and time.

### Application icon

- Custom Consultant Timer clock/timer icon.
- Embedded directly into the Windows executable.
- Used as the `.exe` and desktop icon.
- Used in the application title bar.
- Used in the system tray.

## Technical details

- Standalone **64-bit Windows `.exe`**.
- No installer required.
- No .NET installation required.
- No external database.
- No background cloud service.
- Approximately **2 MB**.
- Local data is stored under `%LOCALAPPDATA%\ConsultantTimer`.
- Designed specifically as a minimal daily consultant work timer rather than a project-management or timesheet platform.

## Windows SmartScreen/code signing

Consultant Timer is currently **not digitally signed with a trusted Windows code-signing certificate**. Because of this, Windows SmartScreen may display an **Unknown publisher** or security warning when launching the application.

This does not indicate that Consultant Timer requires internet access or sends data externally. To provide normal trusted-publisher verification and reduce SmartScreen warnings for distributed builds, the executable should be signed with a trusted Windows code-signing certificate.
