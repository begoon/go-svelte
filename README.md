# Go + Svelte

## Intro

This project uses Svelte (not SvelteKit) and Go to build a single executable web application with modern dynamic web pages powered by Svelte.

### Note - what this project is and what isn't?

This **is a demo project** which can be used as a skeleton for other projects.

This project is **not a library or framework**. The best way is to clone or fork the repo and start tweaking it into your own application. The author is happy to take PRs, but those PR should improve the project demo purpose.

## Why may you try the project?

Go is a perfect tool to build networks and web servers. Try this project if you want to build a web server and host your web app, a self-contained single binary. The single binary can be deployed in a tiny scratch container to GCP or AWS or dropped as an executable to a VPS.

You may start with HTMX if your web app does not require complex reactive behaviour. If it does, you add Alpine JS for more complicated visual effects, but still without a build step for the HTML/JS frontend.

Finally, you should try Svelte if you develop an advanced reactive web application. Svelte is a fantastic framework for building web applications, but it requires the build step and routing on the server if your application has multiple independent pages. The project helps using Go and Svelte together.

## What does the project do?

The Go server provides routing, server-side data loading and injection if necessary. Also, the Go server can do hot reloading during development via a websocket monitor.

The project is for demonstration purposes only. The app exposes a few simple pages with the Counter component in Svelte and links between pages with and without path parameters.

## Important specifics

The hydration is not supported. The Svelte code can use `fetch` directly or `windows.__DATA__`, which is populated by the page handler on the server.

## Required software

- go 1.23+
- Svelte 5+
- bun (to build and package Svelte)
- [air](https://github.com/air-verse/air) for hot reloading
- [just](https://just.systems) as make replacement

## Install dependencies

`just prerequisites`

## Run in dev mode (with hot reloading)

`DEV=1 air`

This command builds the project, Svelte, and Go parts and runs the development server on `:8000`.

When `DEV=1` variable is provided, the Go server runs a /ws websocket endpoint and injects additional javascript to pages, so the browser automatically reloads the page.

If any .go, .svelte, .ts, .js, .css file changes, air rebuilds the application and restarts the server. The web page in the browser reloads automatically.

## Dockerfiles

The application is tested to run on GCP Cloud Run. The application is built into a single binary from "scratch" base image. The root SSL/TLS certificates come from the "golang.org/x/crypto/x509roots/fallback" module, so the Go application can use SSL/TLS.

There are two Dockerfiles provided:

- Dockerfile expects the application to be built locally via `just build`.
- Dockerfile-full is multi-stage. It builds `main.go` itself. However, the Svelte/Vite part must still be built locally via `just build-vite` to produce the `dist/` folder. It is possible to build the Svelte/Vite part in the Dockerfile as well, but it was not the purpose of this project.

## Data loading for Svelte

The Go server injects the page-specific data via the `window.__DATA__` variable, which the Svelte part can use.

## Routing

The `route()` wrapper function in `main.go` provides data injection if necessary. Each page can have specific data. The `path` parameter of `route()` specifies the location of the page content in the `dist/` folder.

A page can be available from multiple routes, specified in multiple `http.HandleFunc("GET...`.

## Svelte page

Each page is an indendent Svelte SPA with index.html, main.ts and the ".svelte" file like "App.svelte".

The pages are configured in `vite.config.ts` in `rollupOptions`.
