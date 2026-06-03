# Noosphere Frontend & Agent Guidelines

This document outlines the coding standards, architecture, and integration patterns for building future features in the Noosphere web application (`apps/web`). AI agents (such as Antigravity) and developers must strictly follow these rules.

---

## 🛠 Tech Stack Overview

- **Core Framework:** Next.js (App Router, Version 16.2.7), React 19 (React Compiler active).
- **Language:** TypeScript (strict mode, path alias `@/*` -> `src/*`).
- **Styling:** Tailwind CSS v4.0 (OKLCH color system, shadcn design tokens).
- **State & Data Fetching:** TanStack Query (React Query v5) + Axios client.
- **Form Handling:** React Hook Form + Zod schema validation.
- **UI Components:** Radix UI primitives + modular Tailwind-styled components.

---

## 📂 Code Architecture & Directory Layout

We enforce a strict **Feature-Sliced Design** pattern under `src/features/`. Any feature (e.g., `chat`, `auth`, `settings`) must follow this internal file layout:

```
src/features/[feature-name]/
├── components/          # Private UI components used only in this feature
│   ├── chat-input.tsx
│   └── chat-window.tsx
├── services/            # API call modules (Axios client wrappers)
│   └── chat-api.ts
├── hooks/               # Feature-specific custom React Query hooks
│   ├── use-chat-messages.ts
│   └── use-send-message.ts
├── types.ts             # Feature-specific TypeScript interfaces & schemas
└── index.ts             # Public API barrel file (exports ONLY public interface elements)
```

### Critical Rules:
1. **No direct imports from private components:** Code outside `src/features/[feature-name]/` must only import from the feature's root `index.ts` (e.g., `import { ChatWindow } from '@/features/chat'`). Never import from internal folders directly (e.g., `import { ChatInput } from '@/features/chat/components/chat-input'`).
2. **Lean App Routing:** Files under `src/app/` (e.g., `src/app/chat/page.tsx`) must be thin wrapper pages. They should merely render components exported from the features layer and carry no business logic or state.
3. **Common Shared UI:** Generic components like buttons, inputs, dialogs, etc., live in `src/components/ui/` (managed via shadcn-style component configurations).

---

## 🎨 Theme, Styling, & UI Conventions

1. **Tailwind CSS v4:** Do not use legacy Tailwind v3 config files. Customize tailwind themes directly in `src/app/globals.css` using the `@theme inline` block and standard CSS variables in OKLCH.
2. **Theme CSS Variables:** Colors are mapped using CSS custom properties (e.g., `var(--background)`, `var(--primary)`). Use semantic utility classes such as `bg-background`, `text-primary-foreground`, `border-border`.
3. **Class Merging:** Always use `cn(...)` from `@/lib/utils` to merge conditional styles and class names safely.
4. **Dark Mode First:** The app wraps layout pages with the `.dark` class by default. Target a modern, premium dark terminal theme:
   - Deep rich backgrounds (`bg-slate-950`, `bg-black`).
   - Subtle borders/separators (`border-slate-800`).
   - Vibrant accents (`text-indigo-400`, `bg-indigo-600`, `shadow-indigo-600/10`).
   - Interactive elements must use scale transitions (`active:scale-95`), hovers, and glassmorphism/blur (`backdrop-blur`).
5. **No Placeholders:** If an image is needed, generate or request a valid illustration asset. Do not use generic placeholders.

---

## 🔄 Data Fetching & Server State

1. **Axios Client:** Always use the shared Axios instance `api` from `@/lib/api-client` for network requests.
2. **React Query Integration:**
   - Every API fetch must be defined as a React Query query or mutation.
   - Encapsulate all queries/mutations into custom hooks inside the feature's `hooks/` folder.
   - Set client defaults (e.g., `refetchOnWindowFocus: false`) in the global `QueryProvider`.
3. **Optimistic Updates:**
   - Implement optimistic UI updates on mutations (like sending messages) to provide high responsiveness.
   - Cancel outstanding queries, preserve the previous cache state for rollback, set the optimistic data, and roll back on error.
   - Example pattern (from `use-send-message.ts`):
     ```typescript
     onMutate: async (newText) => {
       await queryClient.cancelQueries({ queryKey: ['someKey'] });
       const previousData = queryClient.getQueryData(['someKey']);
       queryClient.setQueryData(['someKey'], (old) => [...(old || []), optimisticItem]);
       return { previousData };
     },
     onError: (err, newText, context) => {
       if (context?.previousData) {
         queryClient.setQueryData(['someKey'], context.previousData);
       }
     }
     ```

---

## 📝 Form Validation Guidelines

1. **React Hook Form + Zod:** Use `react-hook-form` coupled with the `@hookform/resolvers/zod` schema validator.
2. **Schema Invariance:** Always define input schemas explicitly using `z.object({...})` and infer the form data type using `z.infer<typeof schema>`.
3. **Disable States:** Ensure submission buttons are disabled while forms are submitting (`disabled={isSubmitting || disabled}`) or when required fields are empty/invalid.

---

## 🔌 Backend API Integration & DTOs

- **Backend Base URL:** `http://localhost:8080` (or `process.env.NEXT_PUBLIC_API_URL`).
- **Endpoints:**
  - `POST /api/v1/chat/message`
  - `GET /api/v1/chat/session/{sessionID}/history`
- **Data Schemas:** Keep JSON field names aligned with the backend Go struct tags (typically snake_case on payloads e.g. `session_id`, `created_at`, but camelCase in TS variables where convenient, though the DTO mapping should directly reflect the API payloads).
- **Error Handling:** Catch and map global client errors through `api-client` interceptors. Always display clear, descriptive error logs and visual recovery controls (e.g., retry buttons) to the user.

---

## 🚀 Branching & Workflow Rules

- **Branch Naming:** All new features must be built on isolated branches starting with `feat/` or `feature/` (e.g., `feat/authentication-frontend`).
- **Main Branch Protection:** Never push directly to `main`. Always coordinate merges via Pull Requests once build and verification checks pass.
