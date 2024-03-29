# state

[![Build Status][build-badge]][build]
![Project Status][project-status]
[![Go Report Card][report-card-badge]][report-card]
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Coverage-39%25-brightgreen.svg?longCache=true&style=flat)</a>
[![Tagged Version][tag-badge]][tag]

[![Project Logo][logo]][logo-large]

*Distributed world and agent state for Bluest Eel models*

## Getting Started

```shell
$ make up
```
```shell
$ make sqlsh
```
```sql
CREATE DATABASE state;
CREATE TABLE state.world (id INT PRIMARY KEY, name VARCHAR);
INSERT INTO state.world VALUES (1, 'test world');
SELECT * FROM state.world;
```

Then connect to a different node and perform the same select query:
```shell
NODE=db3 make sqlsh
```
```sql
SELECT * FROM state.world;
```
```
  id |    name
+----+------------+
   1 | test world
(1 row)

Time: 3.2932ms
```
<!-- Named page links below: /-->

[logo]: https://raw.githubusercontent.com/bluest-eel/branding/master/logo/Logo-v1-x250.png
[logo-large]: https://raw.githubusercontent.com/bluest-eel/branding/master/logo/Logo-v1.png
[build-badge]: https://github.com/bluest-eel/state/workflows/Go/badge.svg
[build]: https://github.com/bluest-eel/state/actions
[report-card-badge]: https://goreportcard.com/badge/bluest-eel/state
[report-card]: https://goreportcard.com/report/bluest-eel/state
[project-status]: https://img.shields.io/badge/project%20status-planning-violet.svg
[tag-badge]: https://img.shields.io/github/tag/bluest-eel/state.svg
[tag]: https://github.com/bluest-eel/state/tags
