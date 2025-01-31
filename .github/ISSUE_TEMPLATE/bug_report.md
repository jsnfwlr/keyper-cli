---
name: Bug report
about: Create a report to help us improve
title: "[BUG] Short summary here"
labels: bug
assignees: jsnfwlr

---

**Describe the bug**
A clear and concise description of what you were trying to do, what you expected to happen, and what actually happened. 

Example: 
> When creating an ed25519-sk key-pair without a Yubikey connected and `--overwrite` set, the application succeeds when it should really show an error message or prompt you to insert a Yubikey before continuing, like it does when `--overwrite` is off.

**To Reproduce**
Provide a step-by-step guide to reproduce the bug, including any and all steps that need to be taken prior to and after running `keyper` in order to trigger the bug.

Example:
> Steps to reproduce the behavior:
> 1. Remove any Yubikey that may be plugged in
> 2. In a terminal type `keyper keys new -t ed25519-sk --overwrite -f ~/.ssh/test_ed25519_sk`
> 3. Follow the on-screen prompts
> 4. Note the command completes without errors and does not request the Yubikey be inserted
> 5. Type `cat ~/.ssh/test_ed25519_sk.pub` and see the public key was created

**Expected behavior**
A clear and concise description of what you expected to happen. 

Example:
> Before prompting you for any additional information, the application should ask you to insert your Yubikey.

**Environment (please complete the following information):**
 - OS: Mac OS | Windows | Linux
 - Version: (if Linux, please provide the Kernel version and version of the distribution which you can get by typing `uname -a`)
 - Installation method: Downloaded binary | Running in Docker | Compile from source

**Additional context**
Add any other information about the problem here.
