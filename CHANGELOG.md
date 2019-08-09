# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.16] - 2019-07-25
### Added
- Support for autoscaling in COMPSs tasks

### Changed
- Documentation updated
- Database and class structs updated

### Fixed
- Errors when retrieving all tasks fixed

## [1.3.9] - 2019-06-25
### Added
- load 'taskqos' elements (predefined SLA templates used to create the SLAs and QoS actions) from _json_ file
- new method added to REST API to upload the 'taskqos' elements: `POST /api/v1/qos/definitions`
- new methods added to REST API to test new functionalities: `GET /api/v1/qos/definitions/{name}`, `GET /api/v1/qos/definitions`
- new method added support COMPSs tasks

### Changed
- Logs removed / updated
- Documentation updated

### Fixed
- swagger.json errors fixed


## [1.2.10] - 2019-06-18
### Added
- new method: `GET /api/v1/docks/tasksqos`
- installation and configuration document for Kubernetes and Openshift

### Changed
- configuration file structure changed to offer support for multiple orchestrators
- adaptation engine updated
- task definition -json- changed
- task QoS updated

### Fixed
- SLALite (Prometheus) `GetValues` function

## [1.1.3] - 2019-06-13
### Added
- TodoTasks file added

### Changed
- logs and exceptions updated
- documentation updated
- configuration updated
- new method implemented: GET /version

### Fixed
- database error when deleting a task
- configuration file error

## [1.1.2] - 2019-06-11
### Added
- Changelog file added to project

### Changed
- logs and exceptions updated
- documentation updated
- configuration initialization improved
- new methods added to the REST API
- adapters (kubernetes, openshift) improved
- automatic SLA generation

### Fixed
- task termination errors
