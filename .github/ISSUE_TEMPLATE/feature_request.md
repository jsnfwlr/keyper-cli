---
name: Feature request
about: Suggest an idea for this project
title: "[FEATURE] Short summary here"
labels: enhancement
assignees: jsnfwlr

---

**If your feature request is related to a problem with existing functionality, please describe the problem**. *If the request is for a new feature unrelated to existing functionality, you can remove this section*
A clear and concise description of what the problem is. 

Example:
> When on-boarding team members, I get them to generate a new SSH key-pair and send me the pubkey and then use keyper-cli to create an account for them, but there is no way to add the pubkey they provide to their account. 

**Describe the solution you'd like**
A clear and concise description of what you want to happen.

Example:
> It would be great if I could use a flag like `--pubkey` to provide the pubkey value as a string or reference a file (keyper-cli should be able to determine whether the value is a pubkey or a file) and/or be asked if I would like to add one or more pubkeys to the account after adding the user and be able to answer with either the pubkey string or the file path to the pubkey.

**Describe alternatives you've considered**
A clear and concise description of any alternative solutions or features you've considered.

Example:
> Currently, without this functionality, I have the following options:
> - I use the UI to add their key to their account
> - get the new member to log in to the UI and add their own key to their account
> - get the new member to install keyper-cli and generate a new key which is added to their account

**Additional context**
Add any other information to support the feature request here.
