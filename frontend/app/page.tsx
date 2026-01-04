import { Activity, ArrowRight, Radar, Server, ShieldCheck, Target } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";

const metrics = [
  { label: "Subdomains found", value: "24", icon: Target },
  { label: "Open services", value: "6", icon: Server },
  { label: "Certs detected", value: "3", icon: ShieldCheck },
  { label: "Latest IPs", value: "5", icon: Radar }
];

const recentFindings = [
  { title: "api.domain.com", detail: "Port 443 (TLS) · Issued by Let's Encrypt" },
  { title: "staging.domain.com", detail: "Port 22 open · SSH host key fingerprint stored" },
  { title: "cdn.domain.com", detail: "Ports 80/443 open · Cloudflare certificate" }
];

export default function Home() {
  return (
    <div className="grid min-h-screen grid-cols-[260px_1fr] text-foreground">
      <aside className="gradient-panel blurred border-r p-6 shadow-lg">
        <div className="flex items-center gap-3">
          <span className="dot h-3 w-3 rounded-full" />
          <div>
            <p className="text-xs uppercase tracking-[0.18em] text-muted-foreground">goscouter</p>
            <p className="text-base font-semibold">Domain Scanner</p>
          </div>
        </div>

        <Separator className="my-6 h-px w-full opacity-50" />

        <div className="space-y-3">
          <Button className="w-full justify-start gap-2" variant="default">
            <Activity size={16} />
            Start Scan
          </Button>
          <Button className="w-full justify-start gap-2" variant="ghost">
            <Server size={16} />
            View Scans
          </Button>
        </div>

        <div className="mt-6 rounded-lg border bg-card/50 p-4">
          <p className="text-sm font-semibold text-foreground">What it does</p>
          <p className="mt-2 text-sm text-muted-foreground leading-relaxed">
            Kick off a domain scan to enumerate subdomains, open services, TLS certificates, backend IPs, and more.
          </p>
        </div>
      </aside>

      <main className="relative overflow-hidden">
        <div className="pointer-events-none absolute inset-0 opacity-40 blur-3xl">
          <div className="absolute -left-32 top-16 h-64 w-64 rounded-full bg-purple-500/30" />
          <div className="absolute right-0 top-0 h-72 w-72 rounded-full bg-cyan-400/20" />
        </div>

        <div className="relative p-10">
          <div className="flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
            <div>
              <p className="text-xs uppercase tracking-[0.18em] text-muted-foreground">Dashboard</p>
              <h1 className="mt-2 text-3xl font-bold">Scout your domain surface</h1>
              <p className="mt-2 max-w-2xl text-sm text-muted-foreground">
                Enter a domain, launch the scan, and review findings as they stream in.
                This is your launchpad to map backend IPs, subdomains, and exposed services.
              </p>
            </div>

            <div className="flex gap-2 rounded-lg border bg-card/70 p-2 shadow-lg backdrop-blur-md">
              <Button size="sm" variant="ghost" className="gap-2">
                <ShieldCheck size={16} />
                TLS
              </Button>
              <Button size="sm" variant="ghost" className="gap-2">
                <Target size={16} />
                Subdomains
              </Button>
              <Button size="sm" variant="ghost" className="gap-2">
                <Server size={16} />
                Services
              </Button>
            </div>
          </div>

          <Card className="mt-6 border-primary/30 bg-gradient-to-br from-card/80 to-card/40 shadow-lg">
            <CardHeader className="gap-4 lg:flex-row lg:items-center lg:justify-between">
              <div>
                <p className="text-xs uppercase tracking-[0.18em] text-muted-foreground">Start Scan</p>
                <CardTitle className="mt-1 text-2xl">Scan a domain</CardTitle>
                <CardDescription className="mt-2">
                  Fire off a discovery pass for subdomains, certificates, and exposed services.
                </CardDescription>
              </div>
              <Badge variant="secondary" className="flex items-center gap-1">
                <Radar size={14} />
                ready
              </Badge>
            </CardHeader>
            <CardContent className="grid gap-3 lg:grid-cols-[1fr_auto] lg:items-center">
              <Input
                placeholder="example.com"
                aria-label="Domain"
                className="h-12 bg-background/80"
              />
              <Button size="lg" className="justify-center gap-2">
                Launch scan <ArrowRight size={16} />
              </Button>
            </CardContent>
          </Card>

          <div className="mt-6 grid gap-4 md:grid-cols-2 xl:grid-cols-4">
            {metrics.map((metric) => (
              <Card key={metric.label} className="border-border/60 bg-card/70">
                <CardContent className="flex items-center gap-3 pt-6">
                  <div className="flex h-10 w-10 items-center justify-center rounded-md bg-accent/40 text-primary">
                    <metric.icon size={18} />
                  </div>
                  <div>
                    <p className="text-2xl font-semibold">{metric.value}</p>
                    <p className="text-sm text-muted-foreground">{metric.label}</p>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          <Card className="mt-6 border-border/60 bg-card/70">
            <CardHeader className="flex flex-row items-center justify-between space-y-0">
              <div>
                <p className="text-xs uppercase tracking-[0.18em] text-muted-foreground">
                  Recent discoveries
                </p>
                <CardTitle className="mt-2 text-xl">Latest scan highlights</CardTitle>
              </div>
              <Button variant="ghost" size="sm" className="gap-2">
                View scans
                <ArrowRight size={14} />
              </Button>
            </CardHeader>
            <CardContent className="space-y-3">
              {recentFindings.map((item) => (
                <div
                  key={item.title}
                  className="flex items-center justify-between rounded-lg border border-border/60 bg-background/60 px-4 py-3"
                >
                  <div>
                    <p className="text-sm font-semibold">{item.title}</p>
                    <p className="text-xs text-muted-foreground">{item.detail}</p>
                  </div>
                  <Badge variant="outline" className="text-primary border-primary/40">
                    new
                  </Badge>
                </div>
              ))}
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  );
}
