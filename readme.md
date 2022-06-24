# RegFileChecker
This is just an idea for now and is far from being feature complete.

---
Currently, the easiest way is to drag a *.reg file to RegFileChecker.exe to see the result.

- Green entries match
- Yellow have other values
- Red does not exist
- Pink has no access rights

![example](https://github.com/spddl/RegFileChecker/blob/main/example.png?raw=true)

#### For color in the Console:
`REG ADD HKEY_CURRENT_USER\Console /v VirtualTerminalLevel /t REG_DWORD /d 0x00000001 /f`
