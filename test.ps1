.\RegFileChecker.exe -NoInfo -Exit @(Get-ChildItem -Path * -Include *.reg | % { $_.Fullname })
Pause
