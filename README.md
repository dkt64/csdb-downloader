# csdb-downloader

**csdb-downloader** is a command-line program to download and organize csdb releases on your storage

**csdb-downloader.json** is a config file for setup (it's created during first run with some default settings):
* *DownloadDirectory* - place where it will download the files (f.e. your USB stick)
* *NoCompoDirectory* - subfolder name for releases which do not belong to any competition
* *LastID* - last ID downloaded, you can adjust it to look deeper in the history
* *Date* - download only releases newer then provided date (download all if empty)
* *Types* - download only listed types of releases (download all if empty)

**Command-line** options:
* *goback* - how many IDs to go back with downloading/checking -> updates config.LastID
* *start* - force ID number to start from -> updates config.LastID
* *date* - download only releases newer then date in form YYYY-MM-DD -> updates config.Date (!)
* *loop* - use this option to make the program loop forever (sleep for 1 minute after checking)
* *all* - use this option to ignore config.Types and download all types of releases

Use -help option to get info in command-line

**Default release types** (demoscene mode):
* C64 Music
* C64 Graphics
* C64 Demo
* C64 One-File Demo
* C64 Intro
* C64 4K Intro
* C64 Crack intro
* C64 Music Collection
* C64 Graphics Collection
* C64 Diskmag
* C64 Charts
* C64 Invitation
* C64 1K Intro 
* C64 Fake Demo

You can adjust the list for your own interests f.e. "C64 Crack", "C64 Tool" or maybe "Other Platform C64 Tool"...

**Notes**
* date (UTC) parameter is only for comparing the dates, program will not start downloading from provided date, ID number is the primary selector
* *LastID* value in config file is increasing after every download

**Example**

If you would like to download all releases starting from ZOO party 2022 let's have a look what ID was first in that day.
It's https://csdb.dk/release/?id=224980

Command below will download all releases released beggining of 2022-10-28:

```csdb-downloader-win64.exe --start=224980 --date=2022-10-28```

**Have a nice day :)**

*DKT*

__Samar Productions / Oct 2022__