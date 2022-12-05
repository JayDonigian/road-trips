# Template for a new entry

#### Replace all `quoted` text with info for the new entry.

---

### Trip Info
| Trip Info        | Fill in the Blanks          |
|------------------|-----------------------------|
| Date (mm-dd)     | `_________________________` |
| Mileage          | `_________________________` |
| Start Short Name | `_________________________` |
| Start Long Name  | `_________________________` |
| Start Emoji      | `_________________________` |
| End Short Name   | `_________________________` |
| End Long Name    | `_________________________` |
| End Emoji        | `_________________________` |

### Expenses
| Expenses           | Delete or Fill in the Blanks |
|--------------------|------------------------------|
| `_Expense_Name_`   | `_Expense Cost_`             |
| `_Expense_Name_`   | `_Expense Cost_`             |
| `_Expense_Name_`   | `_Expense Cost_`             |
| `_Expense_Name_`   | `_Expense Cost_`             |
| `_Expense_Name_`   | `_Expense Cost_`             |

### Locations
* States:
    * `Add info or delete section`
* National Parks:
    * `Add info or delete section`
* Provinces:
    * `Add info or delete section`
* Canadian National Parks:
    * `Add info or delete section`

### Emoji Story `Type out the day's events using only emoji`

### Journal Entry
* `Journal here`

---

#### This form produces the following JSON blob

```json
{
    "name": "{DATE-mm-dd}",
    "date": "2016-{mm-dd}T00:00:00Z",
    "mileage": "{mileage}",
    "start": {
        "emoji": "{start_emoji}",
        "short": "{start_short_name}",
        "long": "{start_long_name}"
    },
    "end": {
        "emoji": "{end_emoji}",
        "short": "{end_short_name}",
        "long": "{end_long_name}"
    },
    "emoji_story": "{emoji_story}",
    "journal_entry": [
        "{journal_entry_line_1}",
        "{journal_entry_line_2}",
        "..."
    ],
    "expenses": [
        {
            "item": "{expense_1_name}",
            "cost": "{expense_1_cost}"
        },
        {
            "item": "{expense_2_name}",
            "cost": "{expense_2_cost}"
        }
    ],
    "states": [
        "{state_1}",
        "{state_2}"
    ],
    "us_parks": [
        "{park_1}",
        "{park_2}"
    ],
    "provinces": [
        "{province_1}",
        "{province_2}"
    ],
    "ca_parks": [
        "{park_1}",
        "{park_2}"
    ]
}
```

