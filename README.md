# Wake-on-LAN Web Interface (Go + Raspberry Pi)

A simple, modern web interface to send Wake-on-LAN (WOL) packets from a Raspberry Pi behind Tailscale.

---

## ⚙️ 1. Configurable Constants

Edit the following line in `main.go` to match your target device:

```go
const targetMAC = "00:11:22:33:44:55"
```

This is the MAC address of the machine you want to wake. Ideally, the target device should have a **static IP lease** in your router to ensure consistency.

---

## 📦 2. Required Dependency: `wakeonlan` CLI

This project currently uses the system’s `wakeonlan` command to send WOL packets.

Install it:

```bash
sudo apt update
sudo apt install wakeonlan
```

Make sure it's available in the system path.

---

## 🔄 3. Evolution: Native Go WOL (No External Command)

You can replace the system command with a pure Go implementation using:

```
github.com/sabhiram/go-wol
```

### Replace `handleWake` with:

```go
import "github.com/sabhiram/go-wol"

func handleWake(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	packet, err := wol.New(targetMAC)
	if err != nil {
		http.Error(w, "Failed to create packet", 500)
		return
	}

	err = packet.Send("255.255.255.255")
	if err != nil {
		http.Error(w, "Failed to send packet", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("WOL packet sent"))
}
```

This removes the need for the `wakeonlan` binary entirely.

---

## 🧰 4. Additional Notes

* The web server runs on port `8080` by default.
* Accessible only within your **Tailscale network**.
* A `setup_wolserver.sh` script is provided to:

  * Build the Go binary
  * Copy assets
  * Register and start a `systemd` service for autostart

After setup:

```bash
http://<your-tailscale-ip>:8080
```

Access from any Tailscale device to trigger the WOL packet.

