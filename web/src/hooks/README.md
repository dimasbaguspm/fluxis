# API Hooks System

This directory contains React hooks for interacting with the Fluxis API, built on top of React Query (TanStack Query) and the auto-generated API client from `src/interfaces/openapi.generated`.

## Architecture

### Base Hooks

All hooks return a **tuple format** for cleaner destructuring:

```ts
[data, error, status, methods]
```

- **`useApiQuery`** - For GET requests (queries)
  - Returns: `[data, error, { isLoading, isFetching, status, ... }, { refetch }]`

- **`useApiMutation`** - For POST, PUT, PATCH, DELETE requests (mutations)
  - Returns: `[data, error, { isPending, isSuccess, status, ... }, { mutate, mutateAsync }]`

- **`useApiInfiniteQuery`** - For paginated queries
  - Returns: `[pages, error, { hasNextPage, isFetchingNextPage, ... }, { fetchNextPage, refetch }]`

### Resource Hooks

Built on top of base hooks, organized by resource:

- **Auth** (`use-auth.ts`) - Login, register, token refresh
- **User** (`use-user.ts`) - Get current user profile
- **Orgs** (`use-orgs.ts`) - Organization CRUD and member management
- **Projects** (`use-projects.ts`) - Project CRUD and visibility
- **Boards** (`use-boards.ts`) - Board CRUD and column management
- **Sprints** (`use-sprints.ts`) - Sprint CRUD and status transitions
- **Tickets** (`use-tickets.ts`) - Ticket CRUD and movement

## Setup

Wrap your app with `QueryProvider`:

```tsx
import { QueryProvider } from "@/providers/query-provider";

export function App() {
  return (
    <QueryProvider>
      <RouterProvider router={router} />
    </QueryProvider>
  );
}
```

## Usage Examples

### Query with Tuple Destructuring

```tsx
import { useGetOrg } from "@/hooks";

export function OrgProfile({ orgId }: { orgId: string }) {
  const [org, err, { isLoading }, { refetch }] = useGetOrg(orgId);

  if (isLoading) return <div>Loading...</div>;
  if (err) return <div>Error: {err.message}</div>;
  if (!org) return <div>Not found</div>;

  return (
    <div>
      <h1>{org.name}</h1>
      <button onClick={() => refetch()}>Refresh</button>
    </div>
  );
}
```

### Mutation with Tuple Destructuring

```tsx
import { useLogin, setApiToken } from "@/hooks";

export function LoginForm() {
  const [result, err, { isPending }, { mutate }] = useLogin();

  const handleSubmit = (email: string, password: string) => {
    mutate(
      { email, password },
      {
        onSuccess: (data) => {
          setApiToken(data.accessToken);
          // Navigate to dashboard
        },
      }
    );
  };

  return (
    <form onSubmit={(e) => {
      e.preventDefault();
      // Extract form values and call handleSubmit
    }}>
      {isPending && <p>Logging in...</p>}
      {err && <p>Error: {err.message}</p>}
      {/* Form fields */}
    </form>
  );
}
```

### List with Pagination

```tsx
import { useListOrgs } from "@/hooks";

export function OrgsList() {
  const [pageNumber, setPageNumber] = useState(1);
  const [orgs, err, { isLoading }, {}] = useListOrgs(
    { pageNumber, pageSize: 10 },
    { enabled: !!pageNumber }
  );

  return (
    <div>
      {isLoading ? (
        <p>Loading...</p>
      ) : (
        orgs?.items?.map((org) => (
          <OrgCard key={org.id} org={org} />
        ))
      )}
      <Pagination
        page={pageNumber}
        totalPages={orgs?.totalPages || 0}
        onPageChange={setPageNumber}
      />
    </div>
  );
}
```

### Infinite Scroll/Pagination

```tsx
import { useListOrgsInfinite } from "@/hooks";

export function OrgsInfiniteList() {
  const [pages, err, { hasNextPage, isFetchingNextPage }, { fetchNextPage }] = useListOrgsInfinite(
    { pageSize: 10 },
    {
      initialPageParam: 1,
      getNextPageParam: (lastPage, _, lastPageParam) =>
        lastPage.totalPages > lastPageParam ? lastPageParam + 1 : undefined,
    }
  );

  return (
    <div>
      {pages?.map((page) =>
        page.items?.map((org) => (
          <OrgCard key={org.id} org={org} />
        ))
      )}
      {hasNextPage && (
        <button onClick={() => fetchNextPage()} disabled={isFetchingNextPage}>
          {isFetchingNextPage ? "Loading..." : "Load More"}
        </button>
      )}
    </div>
  );
}
```

### Dependent Queries

```tsx
import { useGetProject, useListSprints } from "@/hooks";

export function ProjectSprints({ projectId }: { projectId: string }) {
  const [project, err1, projectStatus] = useGetProject(projectId);
  const [sprints, err2, sprintsStatus] = useListSprints(
    { projectId },
    // Only fetch sprints after project is loaded
    { enabled: !!project }
  );

  if (projectStatus.isLoading) return <div>Loading...</div>;
  if (!project) return <div>Not found</div>;
  if (sprintsStatus.isLoading) return <div>Loading sprints...</div>;

  return (
    <div>
      <h2>{project.name}</h2>
      <div>
        {sprints?.items?.map((sprint) => (
          <SprintCard key={sprint.id} sprint={sprint} />
        ))}
      </div>
    </div>
  );
}
```

## API Client Management

### Setting Token

After login, set the access token:

```tsx
import { setApiToken } from "@/hooks";

setApiToken(accessToken);
```

### Clearing Token

On logout:

```tsx
import { clearApiClient } from "@/hooks";

clearApiClient();
```

### Direct Client Access

For advanced use cases:

```tsx
import { getApiClient } from "@/hooks";

const client = getApiClient();
const response = await client.orgs.orgsList();
```

## Tuple Destructuring Pattern

All hooks use consistent tuple format:

```ts
const [data, error, status, methods] = useHook(...)
```

### Status Object Properties

**Query Status:**
- `isLoading` - Query is fetching and there's no data yet
- `isFetching` - Query is currently fetching
- `isSuccess` - Query completed successfully
- `isError` - Query failed with an error
- `status` - "pending" | "error" | "success"
- `fetchStatus` - "idle" | "fetching" | "paused"

**Mutation Status:**
- `isPending` - Mutation is in progress
- `isSuccess` - Mutation succeeded
- `isError` - Mutation failed
- `isIdle` - Mutation hasn't been triggered
- `status` - "idle" | "pending" | "success" | "error"

**Infinite Query Status:**
- All query status properties plus:
- `hasNextPage` - Are there more pages to fetch
- `hasPreviousPage` - Are there previous pages
- `isFetchingNextPage` - Currently fetching next page
- `isFetchingPreviousPage` - Currently fetching previous page

### Methods Object

**Query Methods:**
- `refetch()` - Refetch the query

**Mutation Methods:**
- `mutate(variables)` - Trigger mutation (returns void)
- `mutateAsync(variables)` - Trigger mutation (returns Promise)

**Infinite Query Methods:**
- `fetchNextPage()` - Load next page
- `fetchPreviousPage()` - Load previous page
- `refetch()` - Refetch from start

## Error Handling

Error object contains response error details:

```tsx
const [data, error, { isError }] = useGetOrg(orgId);

if (isError) {
  // error.code - Error code from API
  // error.message - Error message
}
```

## Cache Management

React Query's built-in cache behavior:

- Queries stale after 5 minutes
- Garbage collected after 10 minutes of non-use
- Refetch on window focus disabled
- Automatic retry once on failure

Manually invalidate:

```tsx
import { useQueryClient } from "@tanstack/react-query";

const queryClient = useQueryClient();
queryClient.invalidateQueries({ queryKey: ["orgs"] });
```

## TypeScript Support

All hooks are fully typed from the generated API client:

```tsx
import type { DomainOrganisationModel } from "@interfaces/openapi.generated";

const [org, err, status, methods]: [
  DomainOrganisationModel | undefined,
  any,
  any,
  any
] = useGetOrg(orgId);
```
