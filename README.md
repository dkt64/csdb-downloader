# csdb-downloader
### samar productions 2021

**csdb-downloader** is a commandline tool to download and organize csdb releases on your storage

**csdb-downloader.json** is a config file for setup (it's created during first run with some default settings):
* *DownloadDirectory* - place where it will download the files (f.e. your USB stick)
* *NoCompoDirectory* - folder name for releases which do not belong to any competition
* *LastID* - last ID downloaded, you can adjust it to look deeper in the history
* *Date* - download only releases newer then provided date (download all if empty)
* *Types* - download only listed types of releases

**commandline** options:
* *goback* - how many IDs to go back with downloading/checking -> update of config.LastID
* *startingID* - force ID number to start from -> update of config.LastID
* *date* - download only releases newer then date in form YYYY-MM-DD -> update of config.Date

**default release types** (demoscene mode):
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

you can adjust the list for your own interests f.e. "C64 Crack" or "Tool"

**notes**
* date parameter is only for comparing the dates, the tool will not start downloading from provided date, id number is the primary selector
* the tool downloads only the files which don't exists in the download folder
* if there is new file for download in a release then it's downloaded
* *LastID* value in config file is increasing after every download

have a nice day :)

*dkt/smr*