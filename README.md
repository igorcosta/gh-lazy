# 🚀 Lazy: Your GitHub Project Supercharger

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

Lazy is your secret weapon for turbocharging GitHub CLI commands, project setups at the enterprise level. Say goodbye to tedious manual configurations and hello to lightning-fast, automated awesomeness! We're so lazy that we took advantage of the existing `awesome github cli tool` and beautify with lazyness.

## 🚀 Features That'll Make You Go "Wow!"

- 🏃‍♂️ Sprint through GitHub issue and milestone creation
- 🤖 Automagically set up GitHub Projects (v2)
- 🧙‍♂️ Customize task templates with the power of JSON
- 🔗 Seamlessly integrate with GitHub CLI like a boss

## 🛠️ Requirements

Before you embark on your Lazy journey, make sure you have:

1. Your lazyness!
2. [Homebrew](https://brew.sh/) installed (because we're fancy like that)
3. A GitHub account (you're not living under a rock, are you?)
4. A valid GitHub token with appropriate permissions (we'll show you how)
5. Basic knowledge of JSON (don't worry, it's not rocket science)
6. A burning desire to automate ALL THE THINGS!

## 🏗️ Installation

Let's get this party started:

1. Fire up your terminal (and try not to feel like a hacker)
2. Install the GitHub CLI (if you haven't already):
   ```sh
   brew install gh
   ```
3. Install the Lazy extension (prepare to be amazed):
   ```sh
   gh extension install lazy
   ```
4. Do a little victory dance 🕺💃

## 🎮 Usage

Time to unleash the power of Lazy:

```sh
gh lazy -reponame "your-awesome-username/your-cool-repo" -token "your_super_secret_github_token" -tasks "path/to/your/amazing/tasks.json"
```

### 🎛️ Available Options (Choose Your Destiny)

```
Lazy - The GitHub Milestone and Issue Wizard
Usage: gh lazy [flags]

Flags:
  -reponame string    Your repository's name (e.g., 'cool-dev/awesome-project')
  -tasks string       Path to your magical tasks JSON file
  -tokenfile string   Path to the file containing your GitHub token (default ".token")

Example:
  gh lazy -reponame cool-dev/awesome-project -tasks ./world-domination-plan.json
```

## 🧙‍♂️ How It Works (Warning: Mind-Blowing Content Ahead)

1. Lazy reads your JSON file faster than you can say "automation"
2. It creates milestones and issues in your GitHub repo like a seasoned project manager on steroids
3. A shiny new GitHub Project (v2) materializes out of thin air
4. Issues are automagically added to the project, leaving you more time for coffee breaks

## ⚙️ Configuration (Where the Magic Happens)

Create a JSON file that's so beautiful, it'll bring a tear to your eye:

```json
{
  "milestones": [
    {
      "title": "Phase 1: World Domination",
      "description": "First step towards global supremacy",
      "due_on": "2025-12-31T23:59:59Z"
    }
  ],
  "issues": [
    {
      "title": "Develop Mind Control Ray",
      "body": "Create a device to control the minds of our competitors",
      "milestone": "Phase 1: World Domination",
      "labels": ["top-secret", "high-priority", "needs-coffee"]
    }
  ]
}
```

## 🤝 Contributing (Join the Lazy Revolution)

Want to make Lazy even more awesome? Here's how:

1. Fork the repo (and star it while you're at it)
2. Create a new branch (`git checkout -b feature/mind-blowing-idea`)
3. Commit your changes (`git commit -am 'Add some mind-blowing feature'`)
4. Push to the branch (`git push origin feature/mind-blowing-idea`)
5. Create a new Pull Request and wait for the applause

## 📜 License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details. (TL;DR: Do whatever you want, just don't blame us if your computer gains sentience)

## 🆘 Support

Stuck? Need help? Just want to chat about the meaning of life?

- Open an issue in our GitHub repo

---

Remember: Stay Lazy, Stay Productive! 😴💻