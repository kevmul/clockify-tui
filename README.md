# Clockify Terminal App

A terminal-based time tracking application that integrates with the Clockify API. Track your time entries with a clean, keyboard-driven interface built with Go and Bubble Tea.

## Features

- **Create Time Entries**: Log time with project selection, time ranges, and task descriptions
- **Edit Existing Entries**: Modify any aspect of your time entries
- **Delete Entries**: Remove unwanted time entries
- **View All Entries**: Browse your time entries in an organized list
- **Project Management**: Select from your Clockify projects
- **Flexible Time Input**: Support for various time formats (9a, 9:30a, 3p, 15:30)

## Installation

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Build the application: `go build -o clockify-app`

## Setup

Before using the app, you need to configure your Clockify credentials:

1. Run the application: `./clockify-app`
2. Navigate to the Settings tab
3. Enter your Clockify API key
4. Select your workspace
5. Save the configuration

Your settings are stored locally in a config file for future use.

## Usage

### Navigation

- **Tab**: Switch between Settings, Entries, and Reports views
- **Arrow Keys** or **j/k**: Navigate through lists
- **Enter**: Select items or confirm actions
- **Esc**: Cancel operations or close modals
- **q**: Quit the application

### Creating Time Entries

1. Navigate to the Entries view
2. Press **n** to create a new entry
3. Follow the step-by-step process:
   - **Date Selection**: Choose the date for your time entry
   - **Project Selection**: Select from your available projects
   - **Time Input**: Enter start and end times (e.g., "9a - 5p")
   - **Task Description**: Add a description for your work
   - **Confirmation**: Review and submit your entry

### Editing Time Entries

1. In the Entries view, navigate to the entry you want to edit
2. Press **e** to edit the selected entry
3. Modify any field using the same interface as creating entries
4. Confirm your changes to update the entry

### Deleting Time Entries

1. Navigate to the entry you want to delete
2. Press **d** to delete the selected entry
3. Confirm the deletion when prompted

### Time Format Examples

The app supports flexible time input formats:
- **12-hour format**: `9a`, `9:30a`, `2p`, `2:30p`
- **24-hour format**: `9`, `9:30`, `14`, `14:30`
- **Time ranges**: `9a - 5p`, `9:30a - 5:30p`

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` | Switch between views |
| `↑/↓` or `j/k` | Navigate lists |
| `Enter` | Select/Confirm |
| `Esc` | Cancel/Back |
| `n` | New entry (in Entries view) |
| `e` | Edit entry (in Entries view) |
| `d` | Delete entry (in Entries view) |
| `q` | Quit application |

## Requirements

- Go 1.25.6 or later
- Valid Clockify account and API key
- Terminal with color support (recommended)

## API Integration

The app integrates with the Clockify REST API to:
- Fetch your workspaces and projects
- Create, read, update, and delete time entries
- Sync data in real-time with your Clockify account

All data modifications are immediately reflected in your Clockify web dashboard and mobile app.
