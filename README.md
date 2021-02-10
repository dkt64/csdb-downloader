# csdb-downloader
### samar productions 2021

**csdb-downloader** is a commandline tool to download and organize csdb releases on your storage

**config.json** file for setup (it's created during first run with some default settings):
- *DownloadDirectory* - place where it will download the files (f.e. your USB stick)
- *NoCompoDirectory* - folder name for releases which do not belong to any competition
- *LastID* - last ID downloaded. you can adjust it to look deeper in the history
- *Date* - download only releases newer then provided date (download all if empty)

**commandline** options:
- *goback* - how many IDs go back for updates -> update of config.LastID
- *startingID* - force ID number to start from -> update of config.LastID
- *date* - download only releases newer then date in form YYYY-MM-DD -> update of config.Date

have a nice demo watching :)

*dkt/smr*