# RegFileChecker
[![Downloads][1]][2] [![GitHub stars][3]][4]

[1]: https://img.shields.io/github/downloads/spddl/RegFileChecker/total.svg
[2]: https://github.com/spddl/RegFileChecker/releases "Downloads"

[3]: https://img.shields.io/github/stars/spddl/RegFileChecker.svg
[4]: https://github.com/spddl/RegFileChecker/stargazers "GitHub stars"

---
Currently, the easiest way is to drag a *.reg file to RegFileChecker.exe to see the result.

- Green entries match
- Yellow have other values
- Red does not exist
- Pink has no access rights

![example](https://github.com/spddl/RegFileChecker/blob/main/example.png?raw=true)

With the CLI flags it is also possible to simplify this process
if you have multiple files and want to have only the differences that are not written or are written incorrectly.

```sh
echo Windows Registry Editor Version 5.00 > diff.reg
RegFileChecker.exe -noinfo -nocolor -exit test1.reg >> diff.reg
RegFileChecker.exe -noinfo -nocolor -exit test2.reg >> diff.reg
```