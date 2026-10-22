package factpack

import (
	"path"
	"slices"
	"strings"

	"hastcode.com/site/internal/content"
)

// entry is one file or directory of the repository root or of .github/
// (path "x" or ".github/x").
type entry struct {
	path string
	dir  bool
}

// ciDirs and ciFiles are the CI configuration the Fact Pack reports, by
// path (root level or under .github).
var ciDirs = map[string]bool{".github/workflows": true, ".circleci": true, ".buildkite": true, ".woodpecker": true}

var ciFiles = map[string]bool{
	".gitlab-ci.yml": true, ".travis.yml": true, "Jenkinsfile": true, "azure-pipelines.yml": true,
	".drone.yml": true, "bitbucket-pipelines.yml": true, ".cirrus.yml": true, "appveyor.yml": true,
	".appveyor.yml": true, "cloudbuild.yaml": true, ".woodpecker.yml": true,
}

// manifestNames are the package manifests the Fact Pack reports at the root,
// and manifestExts the manifests named after the project.
var manifestNames = map[string]bool{
	"Cargo.toml": true, "go.mod": true, "package.json": true, "pyproject.toml": true, "setup.py": true,
	"setup.cfg": true, "requirements.txt": true, "Pipfile": true, "pom.xml": true, "build.gradle": true,
	"build.gradle.kts": true, "build.sbt": true, "Gemfile": true, "composer.json": true, "Package.swift": true,
	"pubspec.yaml": true, "mix.exs": true, "CMakeLists.txt": true, "meson.build": true, "deno.json": true,
	"project.clj": true, "deps.edn": true, "stack.yaml": true, "cabal.project": true,
}

var manifestExts = map[string]bool{".csproj": true, ".fsproj": true, ".sln": true, ".gemspec": true, ".nimble": true, ".podspec": true, ".cabal": true, ".opam": true}

// detectFiles reports which well-known files and directories the entries
// hold, as the paths found ("" when absent). Root entries are considered
// before .github ones and both in name order, so the answer does not
// depend on the API's ordering. The README is set from the readme
// endpoint by the caller, which knows GitHub's own choice.
func detectFiles(entries []entry) content.RepoFiles {
	slices.SortFunc(entries, func(a, b entry) int {
		if da, db := strings.Count(a.path, "/"), strings.Count(b.path, "/"); da != db {
			return da - db
		}
		return strings.Compare(a.path, b.path)
	})
	f := content.RepoFiles{CI: []string{}, Manifests: []string{}}
	for _, e := range entries {
		name := path.Base(e.path)
		stem := strings.ToUpper(strings.TrimSuffix(name, path.Ext(name)))
		root := !strings.Contains(e.path, "/")
		if e.dir {
			switch {
			case ciDirs[e.path]:
				f.CI = append(f.CI, e.path)
			case root && f.Docs == "" && (stem == "DOCS" || stem == "DOC" || stem == "DOCUMENTATION"):
				f.Docs = e.path
			case root && f.Examples == "" && (stem == "EXAMPLES" || stem == "EXAMPLE" || stem == "SAMPLES"):
				f.Examples = e.path
			}
			continue
		}
		switch {
		case ciFiles[name]:
			f.CI = append(f.CI, e.path)
		case root && (manifestNames[name] || manifestExts[path.Ext(name)]):
			f.Manifests = append(f.Manifests, e.path)
		case f.License == "" && isLicense(stem):
			f.License = e.path
		case f.Security == "" && stem == "SECURITY":
			f.Security = e.path
		case f.Contributing == "" && stem == "CONTRIBUTING":
			f.Contributing = e.path
		case f.CodeOfConduct == "" && stem == "CODE_OF_CONDUCT":
			f.CodeOfConduct = e.path
		case f.Dockerfile == "" && root && (stem == "DOCKERFILE" || stem == "CONTAINERFILE"):
			f.Dockerfile = e.path
		}
	}
	return f
}

func isLicense(stem string) bool {
	return strings.HasPrefix(stem, "LICENSE") || strings.HasPrefix(stem, "LICENCE") || strings.HasPrefix(stem, "COPYING") || stem == "UNLICENSE"
}
