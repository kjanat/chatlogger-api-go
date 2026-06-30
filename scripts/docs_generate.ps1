#!/bin/pwsh
param (
	[Parameter(HelpMessage = "Generate OpenAPI v3.1 documentation, defaults to v2 which is better supported by swag.")]
	[switch] $v31,
	[Parameter(HelpMessage = "Output path for the generated documentation.")]
	[string] $outPath = "./docs"
	<#,
	[Parameter(HelpMessage = "Instance name for the generated documentation.")]
	[ValidateSet('OpenAPI', 'OpenAPIv3')]
	[string] $instanceName = 'OpenAPIv3',
	[Parameter(HelpMessage = "Package name for the generated documentation.")]
	[ValidateSet('OpenAPI', 'OpenAPI_v3')]
	[string] $packageName = 'OpenAPI'
	#>
)
# This script generates OpenAPI 3.1 documentation from Go annotations using swag v2

function Prep-Env {
	param (
		$outPath,
		$filePath
	)

	# Ensure we're starting with fresh docs
	if (Test-Path "$PSScriptRoot/../$outPath") {
		Write-Host "Cleaning previous documentation..."
		Get-ChildItem -Recurse -Path "$PSScriptRoot/../$outPath/${filePath}_*" | Remove-Item -Force
	}

	# Ensure the output directory exists
	New-Item -ItemType Directory -Path "$PSScriptRoot/../$outPath" -Force | Out-Null

	# Display current swag version
	$swagVersion = $(swag --version)
	Write-Host "Swag version: " $swagVersion
}

function Test-Run {
	param (
		$outPath,
		$filePath,
		$exitCode = $LASTEXITCODE
	)

	if ($LASTEXITCODE -ne 0) {
		Write-Error "Documentation generation failed with code $LASTEXITCODE"
		exit $LASTEXITCODE
	}

	# Verify the documentation files were created
	if (Test-Path "${outPath}/${filePath}_swagger.json") {
		Write-Host "✅ Documentation generated ${filePath} successfully in ${outPath}/" -ForegroundColor Green
	} else {
		Write-Host "❌ Failed to generate documentation files" -ForegroundColor Red
		exit 1
	}
}

# Change directory to project root to ensure proper parsing
Push-Location "$PSScriptRoot/.."
try {
	switch ($v31) {
		$true {
			Write-Host "Generating OpenAPI v3 documentation..."
			$envPath = "OpenAPIv3"
			Prep-Env -outPath $outPath -filePath $envPath
			swag init `
				--output $outPath `
				<# --generalInfo 'cmd/server/main.go' #> `
				--dir 'cmd/server,internal/api,internal/config,internal/domain,internal/handler,internal/hash,internal/jobs,internal/middleware,internal/repository,internal/service,internal/strategy,internal/version' `
				--exclude $outPath `
				--instanceName $envPath --v3.1 --packageName "docs_v3" `
				--parseInternal --parseFuncBody --requiredByDefault `
				--outputTypes 'yaml,json'
				#--state value                             # Set host state for swagger.json
			Test-Run -outPath $outPath -filePath $envPath
		}
		Default {
			Write-Host "Generating OpenAPI v2 documentation..."
			$envPath = "OpenAPI"
			Prep-Env -outPath $outPath -filePath $envPath
			swag init `
				--output $outPath `
				<# --generalInfo 'cmd/server/main.go' #> `
				--dir 'cmd/server,internal/api,internal/config,internal/domain,internal/handler,internal/hash,internal/jobs,internal/middleware,internal/repository,internal/service,internal/strategy,internal/version' `
				--exclude $outPath `
				--instanceName $envPath --packageName "docs" `
				--parseInternal --parseFuncBody --requiredByDefault `
				#--state value                             # Set host state for swagger.json
			Test-Run -outPath $outPath -filePath $envPath
		}
	}
} finally {
	# Return to original directory
	Pop-Location
}
