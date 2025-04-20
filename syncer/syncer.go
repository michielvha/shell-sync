package syncer

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/michielvha/shell-sync/config"
	"github.com/michielvha/shell-sync/filter"
	"github.com/michielvha/shell-sync/filebrowser"
	"github.com/michielvha/shell-sync/history"
)

func SyncLoop(cfg *config.Config, secretFilter *filter.SecretFilter, stopCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	interval := time.Duration(cfg.SyncIntervalSec) * time.Second

	client := filebrowser.NewClient(cfg.Backend.URL, cfg.Backend.Username, cfg.Backend.Password)
	if err := client.Authenticate(); err != nil {
		log.Printf("[ERROR] Filebrowser authentication failed: %v", err)
		return
	}

	for {
		select {
		case <-stopCh:
			log.Println("Sync loop stopped.")
			return
		case <-time.After(interval):
			log.Println("[SYNC] Starting sync cycle...")
			for _, hist := range cfg.HistoryFiles {
				localPath := hist.Path
				remotePath := hist.Path // You may want to namespace this per user/device

				// Download remote file (if exists)
				tmpRemote := localPath + ".remote.tmp"
				err := client.DownloadFile(remotePath, tmpRemote)
				remoteLines := []string{}
				if err == nil {
					remoteLines, _ = history.ReadLines(tmpRemote)
					_ = os.Remove(tmpRemote)
				} else {
					log.Printf("[SYNC] No remote file for %s or failed to download: %v", remotePath, err)
				}

				// Read local file
				localLines, _ := history.ReadLines(localPath)

				// Merge
				merged := history.MergeHistories(localLines, remoteLines)

				// Secret filter
				filtered := []string{}
				for _, line := range merged {
					if secretFilter != nil {
						if filteredLine, blocked := secretFilter.FilterLine(line); blocked && filteredLine == "" {
							continue // Blocked
						} else if blocked {
							filtered = append(filtered, filteredLine)
						} else {
							filtered = append(filtered, filteredLine)
						}
					} else {
						filtered = append(filtered, line)
					}
				}

				// Write merged+filtered to local file
				if err := history.WriteLines(localPath, filtered); err != nil {
					log.Printf("[ERROR] Writing merged history to %s: %v", localPath, err)
				}

				// Upload merged+filtered to remote
				if err := client.UploadFile(remotePath, localPath); err != nil {
					log.Printf("[ERROR] Uploading history %s: %v", remotePath, err)
				}
			}
		}
	}
}
