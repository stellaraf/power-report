<div align="center">
  <br/>
  <img src="https://res.cloudinary.com/stellaraf/image/upload/v1604277355/stellar-logo-gradient.svg" width=300 />
  <br/>
  <h3>Power Report</h3>
  <br/>
  <a href="https://github.com/stellaraf/power-report/actions?query=workflow%3Agoreleaser">
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/stellaraf/power-report/goreleaser?color=9100fa&style=for-the-badge">
  </a>
  <br/>
  This repository contains source code for Stellar's power reporting toolchain.
</div>

## Usage

### Download the latest [release](https://github.com/stellaraf/power-report/releases/latest)

There are multiple builds of the release, for different CPU architectures/platforms:

| Release Suffix |                Platform |
| :------------- | ----------------------: |
| `linux_amd64`  | Linux, Intel or AMD x86 |
| `linux_armv5`  |     Linux, Raspberry Pi |

Right click the one matching your situation, and copy the link. Run the following commands to download and extract:

```shell
wget <release url>
tar xvfz <release file> power-report
```

### Environment Variables

The following environment variables are required:

| Name                      | Description                              | Example                                                              |
| :------------------------ | :--------------------------------------- | :------------------------------------------------------------------- |
| `NMS_URL`                 | LibreNMS URL                             | `"https://nms.example.com"`                                          |
| `NMS_API_KEY`             | LibreNMS API Key                         | `"1234567890"`                                                       |
| `POWER_REPORT_EMAIL_TO`   | Comma-separated list of email recipients | `"alittlemoretopshelf@stellar.tech,ilikethepunishment@stellar.tech"` |
| `POWER_REPORT_EMAIL_FROM` | Sender email address                     | `"Orion Power Reports <power@orion.cloud>"`                          |
| `POWER_REPORT_EMAIL_HOST` | SMTP Host                                | `"smtp.example.com"`                                                 |
| `POWER_REPORT_EMAIL_PORT` | SMTP Port                                | `"25"`                                                               |

### Run the binary

```console
$ sudo ./power-report
```

## Creating a New Release

This project uses [GoReleaser](https://goreleaser.com/) to manage releases. After completing code changes and committing them via Git, be sure to tag the release before pushing:

```
git tag <release>
```

Once a new tag is pushed, GoReleaser will automagically create a new build & release.
