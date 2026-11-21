# Wake-on-LAN Web Interface (Go + Raspberry Pi)

A simple, modern web interface to send Wake-on-LAN (WOL) packets from a Raspberry Pi behind Tailscale.

---

## ⚙️ 1. Configuration

The server reads its settings from the environment:

| Variable | Required | Example | Role |
|---|---|---|---|
| `WOL_MAC` | yes | `00:11:22:33:44:55` | MAC address of the machine to wake |
| `WOL_ADDR` | no | `:8080` (default) | address the web server listens on |

```bash
WOL_MAC=00:11:22:33:44:55 go run main.go
```

The setup script writes `WOL_MAC` into the systemd service, so the address stays out of the code. Ideally, the target device should have a **static IP lease** in your router to ensure consistency.

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

	packet, err := wol.New(targetMAC)  // targetMAC comes from WOL_MAC
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
* A `setup_wolserver.sh` script is provided (`WOL_MAC=… ./setup_wolserver.sh`) to:

  * Build the Go binary
  * Copy assets
  * Register and start a `systemd` service for autostart

After setup:

```bash
http://<your-tailscale-ip>:8080
```

Access from any Tailscale device to trigger the WOL packet.

---

## License

MIT, see [LICENSE](LICENSE).
