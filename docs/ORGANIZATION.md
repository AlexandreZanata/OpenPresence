# Organization Hierarchy

Single polymorphic N-ary tree for private companies and municipalities.

## Node model

```
Tenant (root — company or municipality)
└── OrgNode (configurable type)
    ├── Inherits AttendancePolicy from parent (overridable)
    ├── Managers / approvers per node
    └── Linked GeofenceZones
```

## Private company example

```
Acme Corp (Tenant)
├── HQ São Paulo (DIVISION)
│   ├── Sales (DEPARTMENT)
│   │   ├── Inside Sales (TEAM)
│   │   └── Field Sales (TEAM) ← external route geofence
│   ├── IT (DEPARTMENT)
│   └── HR (DEPARTMENT)
├── Cuiabá Branch (DIVISION)
│   └── Operations (DEPARTMENT)
│       └── Warehouse A (LOCATION)
└── Projects (DIVISION)
    └── Highway BR-163 Site (WORK_SITE) ← temporary geofence
```

## Municipality example

```
City of Example (Tenant)
├── Health Secretariat (DIVISION)
│   ├── Municipal Hospital (LOCATION)
│   │   ├── Nursing (DEPARTMENT) ← 12×36 schedule
│   │   └── Administration (DEPARTMENT) ← 8h/day
│   └── North UBS (LOCATION)
├── Education Secretariat (DIVISION)
│   └── School A (LOCATION)
└── Public Works (DIVISION)
    └── Active construction (WORK_SITE)
```

## RBAC roles

| Role | Scope | Capabilities |
|------|-------|--------------|
| `SUPER_ADMIN` | Tenant | Full control |
| `ORG_ADMIN` | Node + descendants | Employees, policies, geofences |
| `MANAGER` | Immediate node | Approve suspicious punches, reports |
| `HR_ANALYST` | Tenant | Reports, payroll export, manual adjustments |
| `SECURITY_ANALYST` | Tenant | Fraud review, biometric logs |
| `EMPLOYEE` | Self | Punch, own history |
| `AUDITOR` | Tenant read-only | Export, no writes |

## ABAC rules

- **MANAGER** approves only punches for employees in their org subtree.
- **HR_ANALYST** exports only within tenant boundary.
- **EMPLOYEE** accesses only own punch history.

Authorization enforced in Application layer on every request. See [SECURITY.md](SECURITY.md).
