# Inertia.js Integration with Goravel

This document explains how to use Inertia.js with your Goravel application using React and TypeScript.

## Overview

Inertia.js allows you to create single-page applications without building an API. It works by replacing the traditional server-side routing with client-side routing, while still allowing you to write your controllers like you would in a server-rendered application.

## Setup

The following components have been set up in this project:

1. **Inertia Service Provider**: Initializes and registers the Inertia.js service with the application.
2. **Inertia Helper**: Provides helper functions to make it easier to use Inertia.js in your controllers.
3. **React with TypeScript**: Frontend setup using React and TypeScript for building interactive UIs.
4. **Vite**: Used for building and bundling the frontend assets.
5. **Tailwind CSS**: Utility-first CSS framework for styling.

## How to Use

### Rendering Pages

To render an Inertia page from your controller:

```go
import (
    inertiaHelper "players/app/http/inertia"
)

func homeHandler(ctx http.Context) http.Response {
    return inertiaHelper.Render(ctx, "Home", map[string]interface{}{
        "title": "Welcome to my blog",
        "posts": posts,
    })
}
```

### Creating React Components

Create React components in the `resources/js/pages` directory. For example:

```tsx
// resources/js/pages/Home.tsx
import React from 'react';
import { Head } from '@inertiajs/react';

interface HomeProps {
  title: string;
  posts: any[];
}

const Home: React.FC<HomeProps> = ({ title, posts }) => {
  return (
    <>
      <Head>
        <title>{title}</title>
      </Head>
      <div>
        <h1>{title}</h1>
        {posts.map(post => (
          <div key={post.id}>{post.title}</div>
        ))}
      </div>
    </>
  );
};

export default Home;
```

### Building Assets

To build your frontend assets:

1. Install dependencies:
   ```
   npm install
   ```

2. For development:
   ```
   npm run dev
   ```

3. For production:
   ```
   npm run build
   ```

## Middleware

The Inertia middleware is registered globally in the `routes/web.go` file. This middleware handles Inertia requests and responses.

## Shared Data

You can share data globally with all pages using the Inertia service provider:

```go
inertiaManager.Share("user", user)
```

## Links and Forms

Use Inertia.js's Link component and form helpers to navigate between pages without full page reloads:

```tsx
import { Link } from '@inertiajs/react';

<Link href="/users">Users</Link>
```

For forms:

```tsx
import { useForm } from '@inertiajs/react';

const { data, setData, post, processing, errors } = useForm({
  name: '',
  email: '',
});

const submit = (e: React.FormEvent) => {
  e.preventDefault();
  post('/users');
};
```

## Resources

- [Inertia.js Documentation](https://inertiajs.com/)
- [Inertia.js Go Adapter](https://github.com/petaki/inertia-go)
- [React Documentation](https://reactjs.org/)
- [TypeScript Documentation](https://www.typescriptlang.org/)
