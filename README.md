# 🚀 Lazy: Your GitHub Project v2 Supercharger

[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/sindresorhus/awesome)
[![GitHub stars](https://img.shields.io/github/stars/igorcosta/gh-lazy.svg)](https://github.com/igorcosta/gh-lazy/stargazers)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

> Because life's too short for manual setups!

## 🎭 What's in a Name?

**L**ightweight  
**A**utomated  
**Z**ero-effort  
**Y**ielding-results  

## 🌟 Overview

Lazy is your secret weapon for turbocharging GitHub CLI commands and project setups at the enterprise level. Say goodbye to tedious manual configurations and hello to lightning-fast, automated awesomeness! We're so lazy that we took advantage of the existing `awesome github cli tool` and beautified it with laziness.

## 🚀 Features That'll Make You Go "Wow!"

- 🏃‍♂️ Sprint through GitHub issue and milestone creation
- 🤖 Automagically set up GitHub Projects (v2)
- 🧙‍♂️ Customize task templates with the power of JSON
- 🧨 **Nuke GitHub projects and issues with ease**
- 🔗 Seamlessly integrate with GitHub CLI like a boss

## 🛠️ Requirements

Before you embark on your Lazy journey, make sure you have:

1. Your laziness!
2. [Homebrew](https://brew.sh/) installed (because we're fancy like that)
3. A GitHub account (you're not living under a rock, are you?)
4. A valid GitHub token with appropriate permissions (we'll show you how)
5. Basic knowledge of JSON (don't worry, it's not rocket science)
6. A burning desire to automate ALL THE THINGS!

## 🏗️ Installation

Let's get this party started:

1. Fire up your terminal (and try not to feel like a hacker)
2. Install the GitHub CLI (if you haven't already):

   ```bash
   brew install gh
   ```

3. Install the Lazy extension (prepare to be amazed):

   ```bash
   gh extension install lazy
   ```

4. Do a little victory dance 🕺💃

## 🎮 Usage

Time to unleash the power of Lazy:

### Creating Projects, Milestones, and Issues

```bash
gh lazy create --repo "your-awesome-username/your-cool-repo" --tasks "path/to/your/amazing/tasks.json"
```

#### 🎛️ Available Options for `create`

```bash
Lazy - The GitHub Milestone and Issue Wizard
Usage: gh lazy create [flags]

Flags:
  -r, --repo string         Your repository's name (e.g., 'cool-dev/awesome-project')
  -t, --tasks string        Path to your magical tasks JSON file
  -f, --token-file string   Path to the file containing your GitHub token (default ".token")

Example:
  gh lazy create --repo cool-dev/awesome-project --tasks ./world-domination-plan.json
```

### 🧨 Nuking a Project

Delete a GitHub project and optionally all linked issues.

```bash
gh lazy nuke [--projectid <project_id_or_url>] [--all] [--dry-run]
```

- If you provide the `--projectid` (`-p`) flag, the command will delete the specified project.
- If you omit the `--projectid` flag, the tool will:

  1. **List all your available projects** and allow you to select one interactively.
  2. **Ask if you want to perform a dry run first.**
  3. **Ask if you want to delete all associated issues.**

#### 🎛️ Available Options for `nuke`

```bash
Usage: gh lazy nuke [flags]

Flags:
  -p, --projectid string   Project ID or URL to nuke
  -a, --all                Delete all issues linked to the project
      --dry-run            Show what would happen without making changes

Example:
  gh lazy nuke --projectid https://github.com/users/yourusername/projects/1 --all --dry-run
```

**Examples:**

- **Interactive Mode:**

  ```bash
  gh lazy nuke
  ```

  This will prompt you to select a project and configure options interactively.

- **Dry Run Without Deleting Issues:**

  ```bash
  gh lazy nuke --projectid https://github.com/users/yourusername/projects/1 --dry-run
  ```

- **Dry Run With Deleting Issues:**

  ```bash
  gh lazy nuke --projectid 1 --all --dry-run
  ```

- **Actual Deletion:**

  ```bash
  gh lazy nuke --projectid 1 --all
  ```

---

## 🧙‍♂️ How It Works (Warning: Mind-Blowing Content Ahead)

1. Lazy reads your JSON file faster than you can say "automation."
2. It creates milestones and issues in your GitHub repo like a seasoned project manager on steroids.
3. A shiny new GitHub Project (v2) materializes out of thin air.
4. Issues are automagically added to the project, leaving you more time for coffee breaks.
5. **Need to clean up? Use the `nuke` command to delete projects and issues effortlessly.**

---

## 🤝 Contributing (Join the Lazy Revolution)

Want to make Lazy even more awesome? Here's how:

1. Fork the repo (and star it while you're at it).
2. Create a new branch (`git checkout -b feature/mind-blowing-idea`).
3. Commit your changes (`git commit -am 'Add some mind-blowing feature'`).
4. Push to the branch (`git push origin feature/mind-blowing-idea`).
5. Create a new Pull Request and wait for the applause.

---

## 📜 License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details. (TL;DR: Do whatever you want, just don't blame us if your computer gains sentience)

---

## 🆘 Support

Stuck? Need help? Just want to chat about the meaning of life?

- Open an issue in our GitHub repo.

---

Remember: Stay Lazy, Stay Productive! 😴💻
