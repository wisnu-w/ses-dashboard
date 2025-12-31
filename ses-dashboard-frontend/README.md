# SES Dashboard Frontend

A modern, user-friendly React frontend for monitoring AWS SES email events.

## Features

- **Dashboard Overview**: Real-time metrics and analytics
- **Events Management**: Paginated view of all email events
- **Modern UI**: Built with Tailwind CSS and Lucide icons
- **Responsive Design**: Works on desktop and mobile devices
- **Authentication**: JWT-based login system

## Tech Stack

- **React 19** with TypeScript
- **Vite** for fast development and building
- **Tailwind CSS** for styling
- **React Router** for navigation
- **Axios** for API calls
- **Recharts** for data visualization
- **Lucide React** for icons

## Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn

### Installation

1. Clone the repository
2. Navigate to the frontend directory:
   ```bash
   cd ses-dashboard-frontend
   ```

3. Install dependencies:
   ```bash
   npm install
   ```

4. Create environment file:
   ```bash
   cp .env.example .env
   ```

5. Update the API URL in `.env`:
   ```
   VITE_API_URL=http://localhost:8080
   ```

### Development

Start the development server:
```bash
npm run dev
```

The app will be available at `http://localhost:5173`

### Building for Production

```bash
npm run build
```

### Preview Production Build

```bash
npm run preview
```

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── Layout.tsx      # Main layout with sidebar
│   ├── MetricsGrid.tsx # Dashboard metrics cards
│   ├── EventsTable.tsx # Events table with pagination
│   └── Charts.tsx      # Chart components
├── pages/              # Page components
│   ├── LoginPage.tsx
│   ├── DashboardPage.tsx
│   └── EventsPage.tsx
├── services/           # API services
│   └── api.ts
├── types/              # TypeScript type definitions
│   └── api.ts
├── utils/              # Utility functions
└── App.tsx            # Main app component
```

## API Integration

The frontend communicates with the Go backend API. Make sure the backend is running on the configured API URL.

### Authentication

- Login endpoint: `POST /api/login`
- JWT token stored in localStorage
- Automatic token attachment to API requests

### Available Endpoints

- `GET /api/metrics` - Overall metrics
- `GET /api/events` - Paginated events list
- `GET /api/metrics/daily` - Daily metrics
- `GET /api/metrics/monthly` - Monthly metrics
- `GET /api/metrics/hourly` - Hourly metrics

## Demo Credentials

- Username: `admin`
- Password: `password`

## Features Overview

### Dashboard
- Real-time metrics overview
- Interactive charts for daily performance
- Key performance indicators (KPIs)
- Refresh functionality

### Events Page
- Paginated events table
- Sortable columns
- Event type indicators with colors
- Detailed event information

### Responsive Design
- Mobile-first approach
- Collapsible sidebar for mobile
- Optimized layouts for different screen sizes

## Contributing

1. Follow the existing code style
2. Use TypeScript for type safety
3. Test components thoroughly
4. Follow React best practices

## License

This project is part of the SES Dashboard monitoring system.
import reactDom from 'eslint-plugin-react-dom'

export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...
      // Enable lint rules for React
      reactX.configs['recommended-typescript'],
      // Enable lint rules for React DOM
      reactDom.configs.recommended,
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])
```
