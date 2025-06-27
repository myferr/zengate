# DOCKER.md

A guide to building Zengate with Docker instead of downloading GitHub releases.

---

## 1. **Build the Docker Image**

In the root of your Tauri app (where the Dockerfile is located):

```bash
docker build --build-arg TARGETARCH=aarch64 -t zengate-tauri .
```

> Replace `aarch64` with your operating system arch.

> You can tag it anything (`zengate-tauri` is just an example).

---

## 2. **Extract the Built AppImage**

Since the final image is a minimal scratch container, to get the `.AppImage`, run:

```bash
docker create --name temp zengate-tauri true
mkdir -p ./dist
docker cp temp:/Zengate.AppImage ./dist/Zengate.AppImage
docker rm temp
```

then

1. Make it executable. `chmod +x ./dist/Zengate.AppImage`
2. Run it locally. `./dist/Zengate.AppImage`

_Alternatively, a one-liner to create a container, copy, and remove it:_

```bash
docker cp $(docker create zengate-tauri):/Zengate.AppImage ./dist/Zengate.AppImage
```

You can then run the `.AppImage` locally.

---

## 3. **Inspect Inside the Container (Optional)**

To debug or run the build interactively:

```bash
docker run -it --rm zengate-tauri /bin/bash
```

Then:

```bash
cd /app
bunx tauri build
```

---

## 4. **Verify Output**

After building, check your `dist` folder:

```bash
ls dist
# Expected:
# Zengate.AppImage
```

---

## Done!

You now have a portable `.AppImage` built in a reproducible environment.
