package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/kafka"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type DNSController struct {
	RegisteredDNS []api.DNS
	ready         chan bool
	done          chan bool
	subscriber    *kafka.Subscriber
}

func NewDnsController() *DNSController {
	KafkaURL := config.Remotehost + ":9092"
	brokers := []string{KafkaURL}
	group := "dns-controller"
	Controller := &DNSController{
		ready:      make(chan bool),
		done:       make(chan bool),
		subscriber: kafka.NewSubscriber(brokers, group),
	}
	URL := config.GetUrlPrefix() + config.DNSsURL
	var initialDNS []api.DNS
	_ = httputil.Get(URL, &initialDNS, "data")
	log.Info("fetch %d dns from apiserver", len(initialDNS))
	Controller.RegisteredDNS = initialDNS
	return Controller
}

func (s *DNSController) Setup(_ sarama.ConsumerGroupSession) error {
	close(s.ready)
	return nil
}

func (s *DNSController) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (s *DNSController) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if msg.Topic == msg_type.DNSTopic {
			sess.MarkMessage(msg, "")
			s.DNSHandler(msg.Value)
		}
		if msg.Topic == msg_type.ServiceTopic {
			sess.MarkMessage(msg, "")
			s.ServiceHandler(msg.Value)
		}
	}
	return nil
}

func (s *DNSController) DNSHandler(msg []byte) {
	var message msg_type.DNSMsg
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Error("unmarshal msg err: %s", err.Error())
		return
	}
	if message.Opt == msg_type.Delete {
		for index, dns := range s.RegisteredDNS {
			if dns.Name == message.OldDNS.Name {
				s.RegisteredDNS = append(s.RegisteredDNS[:index], s.RegisteredDNS[index+1:]...)
				fileName := fmt.Sprintf("/etc/nginx/conf.d/%s.conf", dns.Host)
				log.Info("delete %s", fileName)
				os.Remove(fileName)
				return
			}
		}
	}
	exist := false
	for index, dns := range s.RegisteredDNS {
		if dns.Name == message.NewDNS.Name {
			exist = true
			s.RegisteredDNS[index] = message.NewDNS
		}
	}
	if !exist {
		s.RegisteredDNS = append(s.RegisteredDNS, message.NewDNS)
	}
	s.WriteDNS()
}

func (s *DNSController) ServiceHandler(msg []byte) {
	var message msg_type.ServiceMsg
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Error("unmarshal msg err: %s", err.Error())
		return
	}
	if message.Opt == msg_type.Delete {
		for index, dns := range s.RegisteredDNS {
			modified := false
			for indexService, path := range dns.Paths {
				log.Info("ServiceName: %s, TargetName: %s", path.ServiceName, message.OldService.Metadata.Name)
				if path.ServiceName == message.OldService.Metadata.Name && path.ServiceNamespace == message.OldService.Metadata.NameSpace {
					s.RegisteredDNS[index].Paths = append(dns.Paths[:indexService], dns.Paths[indexService+1:]...)
					modified = true
				}
			}
			if len(s.RegisteredDNS[index].Paths) == 0 {
				log.Info("delete dns: %s", dns.Name)
				URL := config.GetUrlPrefix() + config.DNSURL
				URL = strings.Replace(URL, config.NamePlaceholder, dns.Name, -1)
				_ = httputil.Delete(URL)
			} else {
				if modified {
					URL := config.GetUrlPrefix() + config.DNSURL
					URL = strings.Replace(URL, config.NamePlaceholder, dns.Name, -1)
					jsonString, _ := json.Marshal(s.RegisteredDNS[index])
					_ = httputil.Put(URL, jsonString)
				}
			}
		}
		return
	}
	for index, dns := range s.RegisteredDNS {
		modified := false
		for serviceIndex, path := range dns.Paths {
			if path.ServiceName == message.NewService.Metadata.Name && path.ServiceNamespace == message.NewService.Metadata.NameSpace {
				s.RegisteredDNS[index].Paths[serviceIndex].ServiceIP = message.NewService.Status.ClusterIP
				modified = true
			}
		}
		if modified {
			URL := config.GetUrlPrefix() + config.DNSURL
			URL = strings.Replace(URL, config.NamePlaceholder, dns.Name, -1)
			jsonString, _ := json.Marshal(s.RegisteredDNS[index])
			_ = httputil.Put(URL, jsonString)
		}
	}
}

func (s *DNSController) WriteDNS() {
	str := "127.0.0.1 localhost\n# The following lines are desirable for IPv6 capable hosts\n::1 ip6-localhost ip6-loopback\nfe00::0 ip6-localnet\nff00::0 ip6-mcastprefix\nff02::1 ip6-allnodes\nff02::2 ip6-allrouters\nff02::3 ip6-allhosts"
	for _, host := range s.RegisteredDNS {
		hostStr := fmt.Sprintf("%s %s\n", config.Remotehost, host.Host)
		str = hostStr + str
	}
	os.WriteFile("/etc/hosts", []byte(str), 0644)
	for _, host := range s.RegisteredDNS {
		NginxStr := "server {\n\tlisten 80;\n"
		hostStr := fmt.Sprintf("\tserver_name %s\n", host.Host)
		NginxStr += hostStr
		for _, path := range host.Paths {
			pathStr := fmt.Sprintf("\tlocation %s {\n", path.Path)
			NginxStr += pathStr
			proxyStr := fmt.Sprintf("\t\tproxy_pass %s:%s\n", path.ServiceIP, path.ServicePort)
			NginxStr += proxyStr
			NginxStr += "\t}\n"
		}
		NginxStr += "}\n"
		fileName := fmt.Sprintf("/etc/nginx/conf.d/%s.conf", host.Host)
		os.WriteFile(fileName, []byte(NginxStr), 0644)
	}
	exec.Command("systemctl", "restart", "nginx").Run()
}

func (s *DNSController) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	topics := []string{msg_type.DNSTopic, msg_type.ServiceTopic}
	s.subscriber.Subscribe(wg, ctx, topics, s)
	<-s.ready
	<-s.done
	cancel()
	wg.Wait()
}
