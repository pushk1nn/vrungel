# Vrungel

Vrungel is a Proof-of-Concept for a Kubernetes-native SIEM, specifically targetting RBAC events. Vrungel uses Discord as a logging endpoint, where users can triage and respond to logs, activating an automated GitOps workflow to push OPA Gatekeeper constraints to a policy-as-code repository.

## Architecture

Vrungel contains 3 controllers:

1. Setup - Submit Discord bot secret and logging channel, policy-as-code repository data, etc.
2. Rule - Submit roles that are "risky" and should be monitored.
3. RoleBindWatcher - Receive data from Rule controller to watch for cluster role binding events.

Once a Discord bot is [created](https://discordpy.readthedocs.io/en/stable/discord.html), the secret is passed as a value within a CR of kind Setup. The Setup reconcile function is to populate the Discord bot session data.

When CR's are submitted with kind Rule, the controller will reconcile by communicating with the RoleBindWatcher controller to update its watchlist of roles. The RoleBindWatcher is different from most controllers in that it does not watch for events of its own resources, rather for Role Binding events. Its reconcile loop will process these events, checking them against the rules sent over by the Rule controller.

If a role binding is corresponding to a Rule submitted by a CR, it will send an embed by calling a handle within the DiscordBotManager. The embed includes a button that will kick off a GitOps workflow to populate a j2 template and push it to the policy-as-code repo specified in the intail Setup CRD.

## Demo

https://youtu.be/GEWlC2b9CcI