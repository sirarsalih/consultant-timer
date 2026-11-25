# Consultant Timer

Consultant Timer is a lightweight, local-first Windows desktop application for simple daily consultant time tracking.

## Features

### Time tracking

- **START** begins tracking working time.
- **START** and **PAUSE** are combined into a single stateful button.
- When the timer is stopped or paused, the button is the green **START** button.
- When tracking is active, the same button becomes the yellow **PAUSE** button.
- Clicking **PAUSE** manually pauses tracking without losing accumulated time.
- **STOP** stops tracking and clears today's accumulated time.
- Time continues accurately across multiple Start/Pause sessions during the same day.
- Saved time survives closing the program and restarting Windows.
- While tracking is active, accumulated time is saved locally so that the latest successfully saved time can be restored after a restart, shutdown, or unexpected power loss.
- The running state is saved locally and restored when the application reopens.
- If the timer was running before Windows restarted, tracking resumes automatically from the saved time.
- If tracking was paused or stopped, it remains inactive after restarting.
- Time spent while Windows is shut down is not counted as working time.

### Two simultaneous time formats

Both formats are displayed at the same time:

- **Decimal hours:** `7.50 hours`
- **Hours/minutes:** `07:30`
- Both values can be selected and copied.

### Automatic idle detection

- Detects **10 minutes of Windows inactivity**.
- Automatically pauses the timer when the idle threshold is reached.
- Removes the detected idle period instead of counting it as working time.
- Displays that the timer was paused because of inactivity.
- Detects new mouse or keyboard activity when the user returns.
- Automatically resumes time tracking after returning from an automatic idle pause.
- The main interface automatically returns to the **Running** state and the button changes to the yellow **PAUSE** button.
- The system-tray state, tooltip, and menu automatically update when tracking resumes.
- A manual **PAUSE** does **not** automatically resume when mouse or keyboard activity is detected; automatic resume only applies to idle-triggered pauses.

### Local and offline operation

- No account required.
- No cloud service.
- No internet connection required.
- Time data is stored locally on the computer.
- Daily data persists between application sessions.
- Previous days remain stored locally when a new day begins.

### Start with Windows

- Consultant Timer automatically starts when the current Windows user signs in.
- No startup checkbox is required.
- Uses the current Windows user's startup configuration.
- Launch the executable once from its intended permanent location to register its startup path.
- Restores the saved timer state and accumulated time after a restart.

### Windows interface

- Native, lightweight Windows desktop application.
- Large and simple daily-time display.
- Displays the current local date in `dd.mm.yyyy` format above the timer, automatically updating when a new day begins.
- Compact window layout without the former startup checkbox or unnecessary space beneath the buttons.
- Combined **START / PAUSE** button:
  - Green **START** while stopped or paused.
  - Yellow **PAUSE** while actively tracking.
- **Red STOP** button.
- The START/PAUSE button automatically changes its text and appearance depending on the current timer state.
- Hand/pointer cursor when hovering over action buttons.
- Displays the current state, such as **Running**, **Paused**, or **Stopped**.
- Shows that idle detection is set to **10 minutes** and that data is stored locally.

### System tray

- Consultant Timer includes a Windows system-tray icon.
- Right-clicking the tray icon provides:
  - **Start**
  - **Pause**
  - **Stop / Clear today**
  - **Show**
  - **Exit**
- Double-clicking the tray icon shows the main window.
- The tray tooltip displays the current timer state, local date in `dd.mm.yyyy` format, decimal hours, and hours/minutes.
- Example: `Consultant Timer - Working - 25.11.2026 - 0.00 h (00:00)`.
- While time tracking is actively **Running**, the system-tray icon displays a prominent **red recording dot** in the lower-right corner.
- The recording dot is restored immediately when the application starts in the Running state after a Windows restart.
- System-tray icon registration is retried if Windows Explorer is not ready at startup or the icon update fails.
- The red recording dot disappears automatically when tracking is paused, stopped, or automatically paused because of inactivity, and returns when tracking resumes.

### Application icon

- Custom Consultant Timer clock/timer icon.
- Embedded directly into the Windows executable.
- Used as the `.exe` and desktop icon.
- Used in the application title bar.
- Used in the Windows taskbar.
- Used in the system tray.
- While the timer is actively **Running**, the Windows taskbar icon displays a **red recording dot in the upper-right corner**.
- While the timer is actively **Running**, the system-tray icon displays a **red recording dot in the lower-right corner**.
- The recording indicators are state-driven: they appear only while working time is actively being recorded and disappear immediately when recording stops or pauses.
- The recording indicators also disappear during an automatic idle pause and return automatically when tracking resumes after mouse or keyboard activity.
- The `.exe` and desktop file icon remain the normal clock/timer icon without the recording dot.

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
