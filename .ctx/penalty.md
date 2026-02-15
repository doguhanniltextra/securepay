# User Caught Errors

- This markdown can only be modified by the user.

## File Permissions

| Filepath      | Permission | Action on Completion              |
|---------------|------------|-----------------------------------|
| [tasks.json](\\wsl.localhost\Ubuntu\home\doguhan\securepay\securepay\tasks.json)  | READ-ONLY  | **REPORT ONLY** (Do NOT edit file)|
| .ctx/*.md   | READ-ONLY | **REPORT ONLY** (Do NOT edit file)                    |

## APPLICATION PROTOCOLS (STRICT)

| Scenario | Forbidden Pattern | Mandatory Pattern | Reason |
| :--- | :--- | :--- | :--- |
| **User Interaction** | Emoji (e.g., ✅, 🚀) | Plain text only | Professional tone requirement. |
| **Task Start** | (Silence) | "I read, starting" | Explicit confirmation requirement. |
| **Shell Commands** | `command1 && command2` | `command1` <br> `command2` | PowerShell compatibility issue. |

### PRE-OPERATION CHECKLIST
Verify the following before performing any operation:
1. [] Am I using emojis? -> STOP. Delete.
2. [] Is this the first response? -> SAY "I read, starting".
3. [] Is there `&&` in the command? -> REWRITE as separate commands.