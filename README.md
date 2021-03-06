# csdb-downloader

**csdb-downloader** is a command-line program to download and organize csdb releases on your storage

**csdb-downloader.json** is a config file for setup (it's created during first run with some default settings):
* *DownloadDirectory* - place where it will download the files (f.e. your USB stick)
* *NoCompoDirectory* - folder name for releases which do not belong to any competition
* *LastID* - last ID downloaded, you can adjust it to look deeper in the history
* *Date* - download only releases newer then provided date (download all if empty)
* *Types* - download only listed types of releases

**Command-line** options:
* *goback* - how many IDs to go back with downloading/checking -> updates config.LastID
* *start* - force ID number to start from -> updates config.LastID
* *date* - download only releases newer then date in form YYYY-MM-DD -> updates config.Date (!)
* *loop* - use this option to make the program loop forever (sleep for 1 minute after checking)

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
* date parameter is only for comparing the dates, program will not start downloading from provided date, ID number is the primary selector
* program downloads only the files which don't exists in the download folder
* if there is new file for download in a release then it will be downloaded
* *LastID* value in config file is increasing after every download

**Example**

If today's 2021-02-11 and newest ID is 199902 (Seraphim by The Solution)
- https://csdb.dk/release/?id=199902

If you would like to go back to the beginning of the year 2021 you should go back with ~1400 releases (198502)
- https://csdb.dk/release/?id=198502

Command below will download all last 1400 releases (IDs) released after 2021-01-01:

```csdb-downloader-win64.exe -goback=1400 --date=2021-01-01```

This one will download all last 1400 releases (IDs) together with all findings from the past:

```csdb-downloader-win64.exe -goback=1400 --date=1980-01-01```

This one will download all releases since 2021-01-01:

```csdb-downloader-win64.exe --start=198535 --date=2021-01-01```

**Have a nice day :)**

*dkt*

__Samar Productions / Feb 2021__