package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/linkerd/linkerd2/controller/gen/controller/discovery"
	"github.com/linkerd/linkerd2/pkg/addr"
	"github.com/spf13/cobra"
)

type endpointsOptions struct {
	namespace    string
	outputFormat string
}

func newEndpointsOptions() *endpointsOptions {
	return &endpointsOptions{
		namespace:    "default",
		outputFormat: "",
	}
}

func newCmdEndpoints() *cobra.Command {
	options := newEndpointsOptions()

	cmd := &cobra.Command{
		Use:   "endpoints [flags]",
		Short: "Display service discovery endpoint state",
		Long:  `Display service discovery endpoint state.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := cliPublicAPIClient()
			endpoints, err := client.Endpoints(context.Background(), &discovery.EndpointsParams{})
			if err != nil {
				return fmt.Errorf("Endpoints API error: %v", err)
			}

			output := renderEndpoints(endpoints, options)
			_, err = fmt.Print(output)

			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&options.namespace, "namespace", "n", options.namespace, "Namespace of the specified resource")
	cmd.PersistentFlags().StringVarP(&options.outputFormat, "output", "o", options.outputFormat, "Output format; currently only \"table\" (default) and \"json\" are supported")

	return cmd
}

func renderEndpoints(endpoints *discovery.EndpointsResponse, options *endpointsOptions) string {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 0, 0, padding, ' ', 0)
	writeEndpointsToBuffer(endpoints, w, options)
	w.Flush()

	return string(buffer.Bytes())
}

type rowEndpoint struct {
	ip       string
	port     uint32
	pod      string
	version  string
	weight   uint32
	identity string
}

func writeEndpointsToBuffer(endpoints *discovery.EndpointsResponse, w *tabwriter.Writer, options *endpointsOptions) {
	rows := []rowEndpoint{}

	for serviceID, servicePort := range endpoints.GetServicePorts() {
		for port, podAddrs := range servicePort.GetPortEndpoints() {
			for _, podAddr := range podAddrs.GetPodAddresses() {
				name := podAddr.GetPod().GetName()
				parts := strings.SplitN(name, "/", 2)
				if len(parts) == 2 {
					name = parts[1]
				}
				row := rowEndpoint{
					ip:       addr.PublicIPToString(podAddr.GetAddr().GetIp()),
					port:     port,
					pod:      name,
					version:  podAddr.GetPod().GetResourceVersion(),
					weight:   1, // TODO: this is hardcoded in the control-plane
					identity: serviceID,
				}

				rows = append(rows, row)
			}
		}
	}

	switch options.outputFormat {
	case "table", "":
		if len(rows) == 0 {
			fmt.Fprintln(os.Stderr, "No endpoints found.")
			os.Exit(0)
		}
		printEndpointTables(rows, w)
	case "json":
		// TODO
		// printEndpointJSON(rows, w)
	}
}

func printEndpointTables(rows []rowEndpoint, w *tabwriter.Writer) {
	headers := []string{
		"IP",
		"PORT",
		"POD",
		"VERSION",
		"WEIGHT",
		"IDENTITY\t", // trailing \t is required to format last column
	}

	fmt.Fprintln(w, strings.Join(headers, "\t"))

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].identity < rows[j].identity
	})

	for _, row := range rows {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%d\t%s\t\n",
			row.ip,
			row.port,
			row.pod,
			row.version,
			row.weight,
			row.identity,
		)
	}
}
