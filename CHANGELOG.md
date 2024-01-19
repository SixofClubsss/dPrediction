# Changelog

This file lists the changes to dPrediction repo with each version.

## 0.3.1 - January 19 2024

### Added
* Scores loading transition 

### Changed
* Go 1.21.5
* Fyne 2.4.3
* dReams 0.11.1
* GetSportsFinals to use Gnomon
* Clean up `rpc` client var names


## 0.3.0 - December 23 2023

### Added

* CHANGELOG
* Pull request and issue templates
* `semver` versioning 
* Stand alone dApp theme
* Asset tabs with profile
* Sync screens

### Changed

* Fyne 2.4.1
* dReams 0.11.0
* Icon resources 
* Confirmations to dialogs 
* implement `gnomes` and funcs
* implement `menu` ShowTxDialog and ShowConfirmDialog
* implement `rpc` PrintError and PrintLog

### Fixed

* Call back to game date in GetWinner funcs
* Deprecated fyne.TextTruncate
* Deprecated container.NewMax
* Fyne error when downloading custom
* Validator hangs