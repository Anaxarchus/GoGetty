# GoGetty

GoGetty is a simplistic and agnostic dependency manager designed for projects that rely on external Git repositories. It manages dependencies effeciently by caching shallow clones, saving space and reducing download sizes. Dependencies are integrated through symbolic links. The tool offers precise version control through the ability to specify branches and commit hashes, allowing you to track exact versions of your dependencies.

While GoGetty was made with Godot in mind, there is nothing about it that prevents it from managing projects on any other platform.

## Building GoGetty

`Go`: GoGetty requires Go version 1.11 or higher. You can download it from [the official Go website](https://golang.org/dl/).
```bash
# Clone the GoGetty repository
git clone https://github.com/Anaxarchus/GoGetty.git

# Navigate to the GoGetty directory
cd GoGetty

# Build the project (Assuming Go is installed)
go build
```

## Installation

There is no need to install GoGetty if you don't want to. Just drop it into your project directory, and optionally add it to your .gitignore:
```bash
cd path/to/your/project
./gogetty init
```

If you would like to install GoGetty, then you can through the install command. The install command will copy GoGetty into the cache location at user/home/.gogetty, and add it to your system PATH.

[ONLY IMPLEMENTED ON LINUX]
```bash
./gogetty install
```

To uninstall:
```bash
gogetty uninstall
```

## Usage

### Initializing a New Project

Start by initializing GoGetty in your project directory:

```bash
cd path/to/your/project
gogetty init
```

This command creates a `.gogetty` configuration file in your project directory. This file is used to track your project's dependencies, and configures the location where module links will be placed. The default location for module links is `modules`. The module links directory is also automatically added to your projects .gitignore, if one is found.

### Adding a Dependency

```bash
cd path/to/your/project
gogetty add <git-repo-url> [--branch <branchName>] [--commit <commitHash>] [--directory <commaSeperatedDirectories>]
```

The branch, commit and directory flags are optional. If empty, the most recent commit of the main branch will be pulled, and the entire repo will be added as a dependency. If the directory flag is passed then only the specified directories will be added to your project as dependencies. This is useful for many Godot respositories since you'll likely want to ignore the root, which typically contains a .project file.

### Updating a Dependency

```bash
cd path/to/your/project
gogetty update <dependencyName> [--branch <branchName>] [--commit <commitHash>] [--directory <commaSeperatedDirectories>]
```

### Removing a Dependency

```bash
cd path/to/your/project
gogetty remove <dependencyName>
```

### Listing Dependencies

```bash
cd path/to/your/project
gogetty list
```

### Cleaning Up Dependencies

```bash
gogetty clean
```

This command will check all registered dependencies and remove any that are no longer valid or needed. If a module's dependent count reaches 0, it will be deleted from the cache.
