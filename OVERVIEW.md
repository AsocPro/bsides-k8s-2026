# Kubernetes: Orchestrating a more secure future

A talk evangelizing kubernetes with a focus on how it's adoption can allow for more secure app deployments.

## How applications have been delivered over time

Brief history of how app deployment has changed, showing the progression toward more secure and reproducible systems.

Methods of Delivery:
- Bare Metal Servers - Direct installation on physical hardware. High performance but difficult to reproduce, long provisioning times, and configuration drift was common. Security patches required careful coordination and often resulted in inconsistent states across servers.
- Virtual Machines - Abstraction layer allowing multiple isolated environments per host. Improved reproducibility through snapshots and templates, but still heavyweight and prone to configuration drift over time. Golden images helped but didn't solve the mutability problem.
- Containers - Lightweight, immutable application packages that bundle code and dependencies. Fast to start, easy to reproduce, and the same image runs identically across environments. Security benefits from reduced attack surface and clear boundaries.

Methods of Deployment:
- Manual Configuration - SSH in, follow a wiki page (if you're lucky). Error-prone, inconsistent, and audit trails are nonexistent. "It works on my server" was a common phrase.
- Bash Scripts - Semi-automated but brittle. Scripts diverge between environments, error handling is often missing, and idempotency is rarely guaranteed.
- Configuration Management (Ansible/Chef/Puppet/Salt/Terraform) - Declarative or semi-declarative configuration as code. Version controlled, auditable, and reproducible. Major step forward but still managing mutable infrastructure.

The Convergence: Following both trends to their logical conclusion, we arrive at immutable deliverables (containers) combined with declarative configuration (infrastructure as code). This is exactly where Kubernetes fits in - it's the orchestration layer that brings these concepts together into a cohesive platform.

## Quick Introduction to Kubernetes

A focused introduction covering the core concepts relevant to security, not a deep-dive tutorial.

What Kubernetes Provides:
- Container Orchestration - Manages the lifecycle of containers across a cluster of machines. Handles scheduling, scaling, and self-healing automatically.
- Declarative Configuration - You define the desired state (YAML manifests), and Kubernetes continuously works to make reality match that state. Changes are version-controlled and auditable.
- Immutable Infrastructure - Instead of patching running containers, you deploy new versions. Rolling updates ensure zero-downtime deployments while maintaining consistency.

Key Concepts for This Talk:
- Pods - The smallest deployable unit, one or more containers sharing network and storage.
- Deployments - Declare how many replicas of a pod should run and how updates roll out.
- Services - Stable networking abstraction that routes traffic to pods.
- Namespaces - Logical isolation boundaries within a cluster for organizing resources and applying policies.

Security Advantages Built-In:
- All configuration changes go through the API server and can be audited.
- Secrets management is a first-class concept (though additional tooling like Vault is often recommended).
- Resource isolation through namespaces, resource quotas, and security contexts.
- The declarative model means your security posture is defined in code and can be reviewed, tested, and version controlled.
- Ease of management as you can manage large number of hardware resources worth in one control plane.

## RBAC: Making sure only the correct people are making the changes

Role-Based Access Control (RBAC) is Kubernetes' authorization mechanism for controlling who can do what within the cluster.

Core RBAC Components:
- Roles/ClusterRoles - Define a set of permissions (verbs like get, list, create, delete) on resources (pods, deployments, secrets, etc.). Roles are namespace-scoped; ClusterRoles are cluster-wide.
- RoleBindings/ClusterRoleBindings - Associate roles with users, groups, or service accounts. This is where you grant the actual permissions.
- Service Accounts - Identity for processes running in pods. Applications should use dedicated service accounts with minimal required permissions.

Security Benefits:
- Principle of Least Privilege - Grant only the permissions needed for each role. Developers might only deploy to their namespace; CI/CD might only update specific deployments.
- Separation of Duties - Different teams get different access levels. Security team can audit, ops can deploy, developers can view logs.
- Auditability - All API requests are tied to an identity. Combined with audit logging, you get a complete record of who did what.

Practical Examples:
- Restrict developers to their team's namespace only
- Allow CI/CD pipelines to update deployments but not delete them
- Give security teams read-only access across all namespaces for auditing
- Prevent anyone from exec'ing into production pods except on-call engineers

## Policy Agents (OPA/Kyverno): Making sure the changes are reasonable

RBAC controls who can make changes, but policy agents control what changes are allowed. They act as admission controllers that validate and mutate resources before they're created.

How They Work:
- Intercept requests to the Kubernetes API server via admission webhooks
- Evaluate resources against defined policies before allowing creation/modification
- Can reject non-compliant resources or mutate them to add required fields

Popular Options:
- OPA/Gatekeeper - Open Policy Agent with Kubernetes-native integration. Policies written in Rego language. Powerful and not limited strictly to kubernetes resources but steeper learning curve.
- Kyverno - Kubernetes-native policy engine. Policies written in YAML, lower barrier to entry. Can validate, mutate, and generate resources.

Security Policies You Can Enforce:
- Require all containers to run as non-root
- Block privileged containers
- Enforce resource limits (CPU/memory) on all pods
- Require specific labels for cost allocation or ownership tracking
- Mandate that images come from approved registries only
- Prevent use of `latest` tag (require explicit versioning)
- Require pods to have security contexts defined

Beyond Admission Control:
- Audit existing resources for compliance violations
- Generate reports on policy violations across the cluster
- Automatically mutate resources to add required fields (e.g., add default network policies, inject sidecars)
- Deploy CA certs to containers automatically.

## Network Policies: Limiting blast radius when something is exploited

By default, all pods in a Kubernetes cluster can communicate with each other. Network policies implement microsegmentation - explicit allow-lists for network traffic between pods.

The Problem They Solve:
- In traditional environments, once an attacker compromises one service, they often have network access to everything on that subnet.
- Without network policies, a compromised pod can reach the database, internal APIs, and potentially pivot to other namespaces.

How Network Policies Work:
- Defined as Kubernetes resources (YAML), applied to pods via label selectors
- Specify allowed ingress (incoming) and egress (outgoing) traffic
- Default-deny policies can block all traffic, then you explicitly allow only required paths
- Requires a CNI plugin that supports network policies (Calico, Cilium, etc.)

Security Benefits:
- Microsegmentation - Each service only communicates with its direct dependencies
- Defense in Depth - Even if application-level controls fail, network-level controls remain
- Reduced Lateral Movement - Attackers can't easily pivot from a compromised frontend to backend databases
- Compliance - Demonstrate network isolation for PCI-DSS, SOC2, etc.

Practical Example:
A typical 3-tier app might have policies like:
- Frontend: Accepts ingress from load balancer, egress only to API tier
- API tier: Accepts ingress from frontend only, egress to database and external APIs
- Database: Accepts ingress only from API tier, no egress allowed

## Looking towards the future

The Kubernetes security ecosystem continues to evolve. Here are emerging technologies that will shape the next generation of secure deployments.

eBPF-Based Security:
- eBPF (extended Berkeley Packet Filter) allows running sandboxed programs in the Linux kernel without modifying kernel source or loading kernel modules.
- Cilium - Uses eBPF for high-performance network policies, going beyond basic L3/L4 filtering to application-aware L7 policies (HTTP, gRPC, Kafka, etc.)
- Tetragon/Falco - Runtime security observability using eBPF. Can detect and alert on suspicious syscalls, file access, network connections in real-time without container modification.
- Talos Linux - Linux based operating system without ssh, traditional userspace or init system. This design lowers the number of potential attack vectors.
- Benefits - Lower overhead than traditional approaches, deeper visibility, and the ability to enforce security at the kernel level before malicious actions complete.

Short-Lived Credentials and Workload Identity:
- Moving away from long-lived static secrets toward dynamic, short-lived credentials.
- Service Account Token Volume Projection - Kubernetes can issue bound, time-limited tokens for pods that automatically rotate.
- Workload Identity Federation - Cloud providers (AWS IRSA, GCP Workload Identity, Azure Workload Identity) allow pods to assume cloud IAM roles without storing cloud credentials in the cluster.
- SPIFFE/SPIRE - Universal workload identity framework. Automatically issues and rotates X.509 certificates for service-to-service authentication.
- Benefits - No more static credentials to leak, automatic rotation, and cryptographic proof of workload identity.

The Bigger Picture:
These technologies represent a shift toward zero-trust architectures where identity is verified continuously, access is granted just-in-time, and security is enforced at every layer from the kernel up. 
