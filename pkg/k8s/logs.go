package k8s

import (
	"bufio"
	"context"

	corev1 "k8s.io/api/core/v1"
)

// WatchLogs watches container's logs of pod in namespace.  It returns channels
// to subscribe log lines.
func (c *Client) WatchLogs(ctx context.Context, namespace, pod, container string) (<-chan string, error) {
	opts := &corev1.PodLogOptions{
		Container: container,
		Follow:    true,
	}
	req := c.clientset.CoreV1().Pods(namespace).GetLogs(pod, opts)
	req.Context(ctx)
	r, err := req.Stream()
	if err != nil {
		return nil, err
	}

	// TODO handle s.Err()
	s := bufio.NewScanner(r)

	ch := make(chan string)
	go func() {
		for s.Scan() {
			ch <- s.Text()
		}
		close(ch)
		defer r.Close()
	}()

	return ch, nil
}
