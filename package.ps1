# 生成压缩包 xx.tar.gz或xx.zip
# 使用 ./package.ps1 -a amd64 -p windows -v v2.0.0

# 任何命令返回非0值退出
$ErrorActionPreference = "Stop"

# 设置默认值
$DefaultOs = $env:OS
$DefaultArch = $env:PROCESSOR_ARCHITECTURE -replace 'x', ''
$PackageDir = ''
$BuildDir = ''
$script:LdFlags = ''
$IncludeFile = @()

# 设置系统、CPU架构
function SetOsArch {
    if (-not $InputOs) {
        $InputOs = $DefaultOs
    }

    if (-not $InputArch) {
        $InputArch = $DefaultArch
    }

    foreach ($Os in $InputOs) {
        if (-not ("linux", "darwin", "windows" -contains $Os)) {
            Write-Host "不支持的系统$Os"
            Exit 1
        }
    }

    foreach ($Arch in $InputArch) {
        if (-not ("386", "amd64" -contains $Arch)) {
            Write-Host "不支持的CPU架构$Arch"
            Exit 1
        }
    }
}

# 初始化
function Init {
    SetOsArch

    if (-not $Version) {
        $Version = Get-GitLatestTag
    }
    $GitCommitId = Get-GitLatestCommit
    $script:LdFlags = "-w -X 'main.AppVersion=$Version' -X 'main.BuildDate=$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')' -X 'main.GitCommit=$GitCommitId'"

    $PackageDir = "${BinaryName}-package"
    $BuildDir = "${BinaryName}-build"

    if (Test-Path $BuildDir) {
        Remove-Item $BuildDir -Force -Recurse
    }
    if (Test-Path $PackageDir) {
        Remove-Item $PackageDir -Force -Recurse
    }

    New-Item -Path $BuildDir -ItemType Directory | Out-Null
    New-Item -Path $PackageDir -ItemType Directory | Out-Null
}

# 获取git最新tag
function Get-GitLatestTag {
    $CommitId = git rev-list --tags --max-count=1
    $TagName = git describe --tags $CommitId

    return $TagName
}

# 获取git最新commit id
function Get-GitLatestCommit {
    return (git rev-parse --short HEAD)
}

# 编译
function Build {
    foreach ($Os in $InputOs) {
        foreach ($Arch in $InputArch) {
            if ($Os -eq "windows") {
                $Filename = "${BinaryName}.exe"
            } else {
                $Filename = $BinaryName
            }
            $Env:CGO_ENABLED = "0"
            $Env:GOOS = $Os
            $Env:GOARCH = $Arch
            go build -ldflags $script:LdFlags -o "${BuildDir}/${BinaryName}-${Os}-${Arch}/${Filename}" $MainFile
        }
    }
}

# 打包
function PackageBinary {
    Set-Location $BuildDir

    foreach ($Os in $InputOs) {
        foreach ($Arch in $InputArch) {
            PackageFile "${BinaryName}-${Os}-${Arch}"
            if ($Os -eq "windows") {
                Compress-Archive -Path "${BinaryName}-${Os}-${Arch}" -DestinationPath "../${PackageDir}/${BinaryName}-${Version}-${Os}-${Arch}.zip" -Force
            } else {
                Compress-Archive -Path "${BinaryName}-${Os}-${Arch}" -DestinationPath "../${PackageDir}/${BinaryName}-${Version}-${Os}-${Arch}.tar.gz" -Force
            }
        }
    }

    Set-Location $PSScriptRoot
}

# 打包文件
function PackageFile {
    if ($IncludeFile.Count -eq 0) {
        return
    }
    foreach ($Item in $IncludeFile) {
        Copy-Item -Path "$PSScriptRoot/../$Item" -Destination $1 -Recurse
    }
}

# 清理
function Clean {
    if (Test-Path $BuildDir) {
        Remove-Item $BuildDir -Force -Recurse
    }
}

# 运行
function Run {
    Init
    Build
    PackageBinary
    Clean
}

function PackageGocron {
    $BinaryName = 'gocron'
    $MainFile = "./cmd/gocron/gocron.go"
    $IncludeFile = @()

    Run
}

function PackageGocronNode {
    $BinaryName = 'gocron-node'
    $MainFile = "./cmd/node/node.go"
    $IncludeFile = @()

    Run
}

# 读取参数
param (
    [string[]]$InputOs,
    [string[]]$InputArch,
    [string]$Version
)

PackageGocron
PackageGocronNode
