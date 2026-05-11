# Aurora Team Structure and Organization

**Document Version:** 1.0  
**Last Updated:** May 11, 2026  
**Owner:** Engineering Manager

---

## Table of Contents

- [Team Philosophy](#team-philosophy)
- [Current Team Structure](#current-team-structure)
- [Roles and Responsibilities](#roles-and-responsibilities)
- [Team Growth Plan](#team-growth-plan)
- [Skills Requirements](#skills-requirements)
- [Communication Protocols](#communication-protocols)
- [Development Workflow](#development-workflow)
- [Meeting Cadence](#meeting-cadence)
- [Team Culture](#team-culture)

---

## Team Philosophy

### Core Values

**1. Ownership and Autonomy**
- Teams own their domains end-to-end
- Individuals empowered to make decisions
- Clear escalation paths for major decisions
- Trust over micromanagement

**2. Collaboration Over Silos**
- Cross-functional teams working together
- Knowledge sharing is expected, not optional
- Pair programming encouraged
- Regular demos and show-and-tell

**3. Quality and Craftsmanship**
- Code quality is everyone's responsibility
- Technical debt addressed proactively
- Testing is part of development, not separate
- Documentation written as features are built

**4. Continuous Learning**
- Learning time built into schedule
- Conference attendance supported
- Internal tech talks encouraged
- Mentorship program for growth

**5. Work-Life Balance**
- Sustainable pace, no heroics
- Async-first communication
- Flexible working hours
- Results over hours logged

### Team Structure Principles

**Small, Focused Teams**
- Teams of 5-9 people (pizza team size)
- Clear mission and scope
- Minimal dependencies on other teams

**Full-Stack Ownership**
- Teams own frontend + backend + infrastructure
- Reduces handoffs and wait time
- Increases sense of ownership

**Embedded Specialists**
- Specialists (DevOps, Security) embedded in teams
- Not separate departments
- Share knowledge, raise overall expertise

---

## Current Team Structure

### Phase 1 Team (Complete)

**Size:** 3 engineers  
**Duration:** 4 weeks  
**Composition:**
- 2 Backend Engineers
- 1 QA Engineer

**Achievements:**
- Core foundation built
- Clean architecture established
- 85% test coverage achieved
- Comprehensive documentation

### Phase 2 Team (June - July 2026)

**Size:** 7 engineers  
**Duration:** 8 weeks

```
Engineering Manager
       │
       ├─── Backend Team Lead (Senior Backend Engineer)
       │    ├─── Backend Engineer (TypeScript)
       │    └─── Backend Engineer (Python)
       │
       ├─── Frontend Lead
       │    ├─── Frontend Engineer (UI)
       │    └─── Frontend Engineer (SDK/Testing)
       │
       ├─── DevOps Engineer (50% time)
       │
       └─── QA Engineer
```

### Phase 3-4 Team (August - November 2026)

**Size:** 10 engineers  
**Duration:** 14 weeks

```
Engineering Manager
       │
       ├─── Backend Team Lead
       │    ├─── Senior Backend Engineer (Auth/Security)
       │    ├─── Backend Engineer (Features)
       │    ├─── Backend Engineer (Performance)
       │    └─── Backend Engineer (GraphQL/WebSockets)
       │
       ├─── Frontend Lead
       │    ├─── Senior Frontend Engineer
       │    └─── Frontend Engineer
       │
       ├─── DevOps/SRE Team
       │    ├─── DevOps Lead
       │    └─── SRE Engineer
       │
       └─── QA Team
            ├─── QA Lead
            └─── QA Engineer
```

### Phase 5+ Team (2027+)

**Size:** 15-20 engineers  
**Multiple specialized teams**

```
VP of Engineering
       │
       ├─── Core Platform Team (5-6)
       │    ├─── Engineering Manager
       │    ├─── Backend Engineers (3)
       │    ├─── DevOps Engineer
       │    └─── QA Engineer
       │
       ├─── Code Generation Team (4-5)
       │    ├─── Tech Lead
       │    ├─── Language Specialists (3)
       │    └─── QA Engineer
       │
       ├─── Enterprise Team (4-5)
       │    ├─── Tech Lead
       │    ├─── Backend Engineers (2)
       │    ├─── Security Engineer
       │    └─── Compliance Specialist
       │
       └─── Developer Experience Team (4-5)
            ├─── Tech Lead
            ├─── Frontend Engineers (2)
            ├─── Technical Writer
            └─── Developer Advocate
```

---

## Roles and Responsibilities

### Engineering Leadership

#### VP of Engineering (Phase 5+)

**Responsibilities:**
- Overall technical vision and strategy
- Team growth and hiring
- Budget management
- Stakeholder communication
- Technical roadmap alignment with business goals

**Time Allocation:**
- 30% - Strategy and planning
- 30% - People management and 1:1s
- 20% - Stakeholder communication
- 10% - Technical architecture review
- 10% - Recruiting and hiring

**Required Skills:**
- 10+ years engineering experience
- 5+ years in leadership roles
- Strong architectural background
- Excellent communication skills
- Experience scaling teams

**Reports to:** CTO or CEO

---

#### Engineering Manager

**Responsibilities:**
- Day-to-day team management
- Sprint planning and execution
- Remove blockers for team
- Career development and mentorship
- Performance reviews and feedback
- Hiring and onboarding
- Technical decisions with team leads

**Time Allocation:**
- 40% - People management (1:1s, feedback, career growth)
- 30% - Project management (planning, tracking, unblocking)
- 20% - Technical work (architecture, code review)
- 10% - Recruiting and hiring

**Required Skills:**
- 8+ years engineering experience
- 3+ years team lead/manager experience
- Strong technical background
- People management skills
- Excellent communication

**Reports to:** VP of Engineering (or CTO in early phases)

**Phase 1-2:** Manages 3-7 people  
**Phase 3-4:** Manages 10 people  
**Phase 5+:** Manages 3-4 team leads (15-20 total)

---

### Backend Engineering

#### Senior Backend Engineer / Tech Lead

**Responsibilities:**
- Technical leadership for backend team
- Architecture decisions and documentation
- Code review and quality standards
- Mentoring junior engineers
- Performance and scalability
- Production operations
- Hands-on development (60-70% time)

**Time Allocation:**
- 60% - Hands-on development
- 20% - Code review and mentorship
- 10% - Architecture and design
- 10% - Team coordination and planning

**Required Skills:**
- 7+ years backend development
- Expert in Go (or primary language)
- Strong architectural skills
- Database design and optimization
- System design and scalability
- API design expertise
- Production operations experience
- Mentorship and teaching ability

**Key Technologies:**
- Go, PostgreSQL, GORM
- RESTful API design
- Clean Architecture patterns
- Docker, Kubernetes
- Git, CI/CD

**Reports to:** Engineering Manager

**Career Path:**
- Senior Backend Engineer → Staff Engineer → Principal Engineer
- OR Senior Backend Engineer → Tech Lead → Engineering Manager

---

#### Backend Engineer

**Responsibilities:**
- Feature development
- Bug fixing and maintenance
- Writing tests and documentation
- Code review participation
- Operational support

**Time Allocation:**
- 70% - Feature development
- 15% - Bug fixing and maintenance
- 10% - Code review
- 5% - Documentation

**Required Skills:**
- 3-7 years backend development
- Proficient in Go
- Database knowledge (PostgreSQL)
- API design
- Testing practices
- Git workflows

**Key Technologies:**
- Go, PostgreSQL, GORM
- Gin framework
- Docker
- Git, GitHub

**Reports to:** Engineering Manager or Tech Lead

**Career Path:**
- Junior Backend Engineer → Backend Engineer → Senior Backend Engineer

---

#### Language Specialist (Phase 2+)

**TypeScript Specialist:**
- Expert in TypeScript ecosystem
- Node.js, Express.js, NestJS
- Prisma, TypeORM
- Jest, testing frameworks
- npm, package management

**Python Specialist:**
- Expert in Python ecosystem
- FastAPI, Flask, Django
- SQLAlchemy, async Python
- Pytest, testing frameworks
- pip, Poetry, package management

**Go Specialist (Future):**
- Expert in Go ecosystem
- Gin, Echo frameworks
- GORM, database/sql
- Go testing
- Go modules

**Responsibilities:**
- Design and maintain language-specific generators
- Stay current with language best practices
- Template quality and optimization
- Language-specific testing
- Documentation and examples

**Time Allocation:**
- 80% - Generator development
- 10% - Research and learning
- 10% - Documentation

**Required Skills:**
- 5+ years in specific language
- Deep ecosystem knowledge
- Code generation experience (helpful)
- Template design
- Best practices expertise

**Reports to:** Tech Lead (Code Generation Team)

---

### Frontend Engineering

#### Frontend Lead

**Responsibilities:**
- Frontend architecture and decisions
- Component library design
- Performance optimization
- Mentoring frontend engineers
- UI/UX collaboration
- Hands-on development (60-70%)

**Time Allocation:**
- 60% - Hands-on development
- 20% - Architecture and design
- 10% - Mentorship and code review
- 10% - Coordination and planning

**Required Skills:**
- 6+ years frontend development
- Expert in React (or primary framework)
- TypeScript proficiency
- Component library design
- Performance optimization
- Accessibility (a11y)
- Design systems

**Key Technologies:**
- React, TypeScript
- Next.js or Vite
- CSS-in-JS or Tailwind
- Jest, Testing Library
- Webpack or Vite

**Reports to:** Engineering Manager

---

#### Frontend Engineer

**Responsibilities:**
- UI component development
- Integration with backend APIs
- Testing and documentation
- Performance optimization
- Accessibility implementation

**Time Allocation:**
- 70% - Component development
- 15% - Integration and testing
- 10% - Code review
- 5% - Documentation

**Required Skills:**
- 3-6 years frontend development
- Proficient in React
- TypeScript knowledge
- CSS and responsive design
- Testing practices
- Git workflows

**Key Technologies:**
- React, TypeScript
- CSS, Tailwind (or similar)
- REST API integration
- Jest, Testing Library

**Reports to:** Frontend Lead or Engineering Manager

---

### DevOps / SRE

#### DevOps Lead / SRE

**Responsibilities:**
- Infrastructure architecture
- CI/CD pipeline design and maintenance
- Production monitoring and alerting
- Incident response
- Security and compliance
- Cost optimization
- Disaster recovery planning

**Time Allocation:**
- 40% - Infrastructure development
- 30% - Production operations
- 20% - Automation and tooling
- 10% - Team support

**Required Skills:**
- 6+ years DevOps/SRE experience
- Kubernetes expertise
- Cloud platforms (AWS, GCP, or Azure)
- CI/CD tools (GitHub Actions, GitLab CI)
- Monitoring tools (Prometheus, Grafana)
- Infrastructure as Code (Terraform)
- Security best practices
- Incident management

**Key Technologies:**
- Kubernetes, Docker
- AWS or GCP
- Terraform or Pulumi
- GitHub Actions
- Prometheus, Grafana
- PostgreSQL administration

**Reports to:** Engineering Manager

---

#### DevOps Engineer

**Responsibilities:**
- CI/CD maintenance
- Environment management
- Deployment automation
- Monitoring setup
- Infrastructure support

**Time Allocation:**
- 50% - CI/CD and automation
- 30% - Environment management
- 20% - Team support

**Required Skills:**
- 3-5 years DevOps experience
- Docker and Kubernetes
- CI/CD experience
- Cloud platforms
- Scripting (Bash, Python)
- Monitoring tools

**Key Technologies:**
- Docker, Kubernetes
- GitHub Actions or GitLab CI
- AWS or GCP
- Bash, Python scripting
- Prometheus basics

**Reports to:** DevOps Lead or Engineering Manager

---

### Quality Assurance

#### QA Lead (Phase 3+)

**Responsibilities:**
- QA strategy and process
- Test plan creation
- Test automation framework
- Quality metrics and reporting
- Team mentorship
- Production quality monitoring

**Time Allocation:**
- 40% - Test planning and strategy
- 30% - Hands-on testing
- 20% - Automation and tooling
- 10% - Reporting and metrics

**Required Skills:**
- 6+ years QA experience
- Test automation expertise
- API testing tools
- Performance testing
- Security testing basics
- CI/CD integration
- Programming skills (Go, Python, or JavaScript)

**Key Technologies:**
- Go testing framework
- Postman or similar API tools
- Selenium or Cypress
- JMeter or k6 (performance)
- Git, CI/CD

**Reports to:** Engineering Manager

---

#### QA Engineer

**Responsibilities:**
- Test case creation and execution
- Manual and automated testing
- Bug reporting and tracking
- Regression testing
- Documentation testing
- User acceptance testing coordination

**Time Allocation:**
- 50% - Manual testing
- 30% - Test automation
- 15% - Bug reporting and tracking
- 5% - Documentation

**Required Skills:**
- 2-5 years QA experience
- Manual testing expertise
- Basic automation skills
- API testing
- Bug tracking tools
- SQL basics
- Attention to detail

**Key Technologies:**
- Postman or Insomnia
- Browser DevTools
- Basic scripting
- Git basics
- Jira or similar

**Reports to:** QA Lead or Engineering Manager

---

### Supporting Roles

#### Technical Writer (Part-time to Full-time)

**Responsibilities:**
- User documentation
- API reference documentation
- Tutorial creation
- Blog posts and guides
- Documentation maintenance
- Video script writing

**Phase 1-2:** 20% time  
**Phase 3-4:** 50% time  
**Phase 5+:** 100% time

**Required Skills:**
- 3+ years technical writing
- Software development understanding
- Markdown and Git
- API documentation experience
- Clear, concise writing
- Developer empathy

---

#### Product Manager (Part-time initially)

**Responsibilities:**
- Product vision and strategy
- Feature prioritization
- User research and feedback
- Roadmap management
- Stakeholder communication
- Success metrics definition

**Phase 1-2:** 20% time  
**Phase 3-4:** 50% time  
**Phase 5+:** 100% time

**Required Skills:**
- 4+ years product management
- Developer tools experience
- Agile/Scrum knowledge
- Data-driven decision making
- Excellent communication
- Technical background

---

#### Designer (Part-time)

**Responsibilities:**
- UI/UX design
- Design system creation
- User flow design
- Prototyping
- User research support
- Brand consistency

**Phase 1-2:** 10% time  
**Phase 3-4:** 25% time  
**Phase 5+:** 50% time

**Required Skills:**
- 3+ years UI/UX design
- Figma or similar tools
- Design systems
- User research
- Accessibility knowledge
- Developer tools experience

---

## Team Growth Plan

### Phase 1: Foundation (Complete)
**Team Size:** 3  
**Composition:** 2 Backend, 1 QA

### Phase 2: Code Generation (June-July 2026)
**Team Size:** 7  
**Composition:** 3 Backend, 2 Frontend, 1 DevOps (50%), 1 QA  
**New Hires:** +4 (2 Frontend, 1 Backend specialist, upgrade DevOps)

### Phase 3-4: Advanced Features & Operations (August-November 2026)
**Team Size:** 10  
**Composition:** 4 Backend, 2 Frontend, 2 DevOps/SRE, 2 QA  
**New Hires:** +3 (1 Backend, 1 DevOps, 1 QA)

### Phase 5: Enterprise Features (December 2026 - January 2027)
**Team Size:** 12  
**Composition:** 5 Backend (1 Security), 2 Frontend, 2 DevOps, 2 QA, 1 Technical Writer  
**New Hires:** +2 (1 Security Engineer, 1 Technical Writer full-time)

### Phase 6+: Ecosystem & Scale (February 2027+)
**Team Size:** 15-20  
**Composition:** Multiple specialized teams  
**New Hires:** +3-8 (various specializations)

### Hiring Timeline

```
2026:
Q2 (Phase 2)
  - April: Hire 2 Frontend Engineers
  - May: Hire TypeScript Specialist Backend Engineer

Q3 (Phase 3)
  - June: Hire 1 Senior Backend Engineer (Auth/Security)
  - July: Hire 1 DevOps/SRE Engineer
  - August: Hire 1 QA Engineer

Q4 (Phase 4-5)
  - September: Hire 1 Security Engineer
  - October: Hire 1 Backend Engineer
  - November: Technical Writer to full-time

2027:
Q1-Q2 (Phase 6)
  - January: Hire Developer Advocate
  - February: Hire 2 Backend Engineers
  - March: Hire 1 Frontend Engineer
  - April: Hire 1 QA Engineer
```

### Onboarding Timeline

**Week 1: Orientation**
- Company overview and values
- Team introductions
- Development environment setup
- Codebase overview
- Architecture walkthrough

**Week 2-3: Learning Phase**
- Pair programming with team members
- Documentation deep dive
- Small bug fixes or improvements
- Code review participation

**Week 4: First Real Task**
- Assigned to real feature or bug
- Supported by mentor
- First PR submitted

**By Week 6:**
- Fully productive
- Contributing independently
- Comfortable with codebase and process

---

## Skills Requirements

### Core Technical Skills (All Engineers)

**Must Have:**
- Git and GitHub workflows
- Test-driven development
- Code review practices
- CI/CD basics
- Agile/Scrum methodology
- Documentation writing
- Problem-solving mindset

**Nice to Have:**
- Open source contribution experience
- Previous startup experience
- Public speaking or writing
- Cross-functional collaboration

### Backend Engineer Skill Matrix

| Skill | Junior | Mid | Senior | Staff |
|-------|--------|-----|--------|-------|
| **Go** | Basic syntax | Idiomatic Go | Expert, design patterns | Community contributor |
| **Databases** | SQL basics | PostgreSQL, schema design | Optimization, replication | Database architecture |
| **APIs** | REST basics | REST + GraphQL | API design patterns | Industry standards setter |
| **Architecture** | MVC | Clean Architecture | System design | Distributed systems |
| **Testing** | Unit tests | Integration tests | E2E, performance | Testing strategy |
| **DevOps** | Docker basics | Docker + K8s basics | CI/CD, monitoring | Infrastructure expert |
| **Leadership** | - | Code review | Mentoring | Tech direction |

### Frontend Engineer Skill Matrix

| Skill | Junior | Mid | Senior | Staff |
|-------|--------|-----|--------|-------|
| **React** | Basic components | Hooks, context | Performance, patterns | Framework expert |
| **TypeScript** | Basic types | Generics, utilities | Advanced patterns | TS architecture |
| **CSS** | Basics | Responsive, Flexbox | Design systems | CSS architecture |
| **Testing** | Basic tests | Component tests | E2E, visual tests | Testing strategy |
| **Performance** | - | Basic optimization | Profiling, optimization | Performance expert |
| **A11y** | - | Basic a11y | WCAG compliance | A11y champion |
| **Leadership** | - | Code review | Mentoring | Tech direction |

---

## Communication Protocols

### Synchronous Communication

**Daily Standup**
- Time: 9:00 AM daily
- Duration: 15 minutes max
- Format: What I did, what I'll do, blockers
- Tool: Video call (all remote) or in-person
- Who: Full team

**Weekly Planning**
- Time: Monday 10:00 AM
- Duration: 1 hour
- Agenda: Sprint planning, task breakdown, estimation
- Tool: Video call + Jira
- Who: Full team

**Weekly Demo**
- Time: Friday 3:00 PM
- Duration: 30 minutes
- Format: Live demos of completed work
- Tool: Video call + screen share
- Who: Team + stakeholders

**Ad-hoc Pairing**
- As needed
- Screen sharing + video
- Encouraged for complex tasks or learning

### Asynchronous Communication

**Primary: GitHub**
- Pull requests for code review
- Issues for bug tracking
- Discussions for proposals
- Wiki for documentation

**Secondary: Slack**
- `#aurora-general` - Team updates
- `#aurora-dev` - Development discussions
- `#aurora-blockers` - Urgent issues only
- `#aurora-random` - Non-work chat

**Email**
- For external communication
- Formal announcements
- Performance reviews (confidential)

### Communication Guidelines

**Response Time Expectations:**
- Urgent (blocker): < 1 hour
- High priority: < 4 hours
- Normal: < 1 day
- Low priority: < 2 days

**Working Hours:**
- Core hours: 10 AM - 3 PM (team timezone)
- Flexible start/end times
- Respect "Do Not Disturb" status
- No expectation to respond outside working hours

**Code Review SLA:**
- Review requested: Review within 24 hours
- Changes requested: Author responds within 24 hours
- Approved: Merge immediately or within 4 hours

**Documentation:**
- Code: Inline comments for complex logic
- Architecture: ADRs (Architecture Decision Records)
- API: OpenAPI/Swagger specs
- Process: Team wiki or Notion

---

## Development Workflow

### Sprint Cadence

**2-Week Sprints**

**Sprint Planning (Monday Week 1):**
- Review roadmap and priorities
- Break down user stories
- Estimate effort (story points)
- Assign tasks
- Set sprint goal

**Daily Standups (Every day):**
- 15 minutes
- Share progress and blockers
- Adjust plan if needed

**Sprint Review (Friday Week 2):**
- Demo completed work
- Collect feedback
- Update roadmap

**Sprint Retrospective (Friday Week 2):**
- What went well
- What could improve
- Action items for next sprint

### Task Management

**Tool:** Jira (or Linear)

**Board Columns:**
- Backlog
- Todo
- In Progress
- Code Review
- Testing
- Done

**Task Types:**
- Epic (large feature)
- Story (user story)
- Task (technical task)
- Bug (defect)
- Spike (research)

**Story Points:**
- 1 point = 2-4 hours
- 2 points = half day
- 3 points = 1 day
- 5 points = 2-3 days
- 8 points = 4-5 days (consider splitting)
- 13 points = too big, must split

### Git Workflow

**Branching Strategy:**
```
main (production)
  ↑
develop (integration)
  ↑
feature/ABC-123-short-description
bugfix/ABC-456-short-description
hotfix/critical-bug-fix
```

**Branch Naming:**
- `feature/ABC-123-add-code-generation`
- `bugfix/ABC-456-fix-schema-validation`
- `hotfix/emergency-security-fix`

**Commit Messages:**
```
type(scope): short description

Longer description if needed.

Closes #123
```

Types: feat, fix, docs, style, refactor, test, chore

**Pull Request Process:**
1. Create PR from feature branch to develop
2. Fill out PR template
3. Request review from 1-2 people
4. Address feedback
5. Get approval
6. Squash and merge

**Code Review Checklist:**
- [ ] Code works and meets requirements
- [ ] Tests added and passing
- [ ] Documentation updated
- [ ] No security issues
- [ ] Performance is acceptable
- [ ] Code style consistent
- [ ] PR description clear

### Deployment Process

**Environments:**
- **Local:** Developer machine
- **Development:** Auto-deploy from develop branch
- **Staging:** Manual deploy from develop (before release)
- **Production:** Manual deploy from main (after approval)

**Deployment Steps:**
1. Code merged to develop
2. Auto-deploy to development environment
3. Automated tests run
4. Manual deploy to staging for testing
5. QA approval
6. Create release PR (develop → main)
7. Final review and approval
8. Merge and deploy to production
9. Monitor metrics and alerts

---

## Meeting Cadence

### Daily

**Daily Standup (15 min)**
- 9:00 AM
- Full team
- Quick status update

### Weekly

**Sprint Planning (1 hour)**
- Monday 10:00 AM
- Plan week's work
- Full team

**Sprint Review / Demo (30 min)**
- Friday 3:00 PM
- Show completed work
- Team + stakeholders

**Tech Sync (30 min)**
- Wednesday 2:00 PM
- Technical discussions
- Architecture decisions
- Engineering team only

### Bi-Weekly

**Sprint Retrospective (1 hour)**
- End of sprint (Friday)
- Continuous improvement
- Full team

**1:1s with Manager (30 min)**
- Scheduled individually
- Career development
- Feedback
- Private concerns

### Monthly

**All-Hands (1 hour)**
- First Monday of month
- Company updates
- Roadmap review
- Q&A

**Tech Talks (1 hour)**
- Third Thursday
- Learn from each other
- External speakers
- Optional but encouraged

### Quarterly

**Planning Session (4 hours)**
- First week of quarter
- Roadmap for next 3 months
- OKR setting
- Team + leadership

**Team Offsite (Full day)**
- Last month of quarter
- Team building
- Strategic planning
- Fun activities

---

## Team Culture

### Values in Action

**Ownership**
- Engineers own their features end-to-end
- "You build it, you run it" philosophy
- Take initiative, don't wait for permission
- Fix things you notice are broken

**Collaboration**
- Pair programming encouraged
- Ask for help, offer help
- Knowledge sharing is valued
- No knowledge hoarding

**Quality**
- Code review is thoughtful, not rubber-stamping
- Test coverage matters
- Documentation is part of the job
- Technical debt addressed regularly

**Learning**
- 10% time for learning
- Conference budgets ($2000/year per person)
- Internal tech talks
- Experiment with new technologies

**Balance**
- No late-night deploys (except emergencies)
- Sustainable pace
- Take vacation
- Mental health matters

### Recognition and Rewards

**Peer Recognition:**
- Kudos in team meetings
- Slack shout-outs
- Monthly team awards

**Performance Reviews:**
- Quarterly informal check-ins
- Annual formal review
- 360-degree feedback
- Career development plan

**Career Growth:**
- Clear career ladder
- Promotion criteria transparent
- Mentorship program
- Leadership opportunities

**Compensation:**
- Competitive salaries
- Annual raises (merit + market)
- Equity/stock options
- Performance bonuses

---

## Appendix: Job Descriptions

### Senior Backend Engineer - Template Engine Lead

**Location:** Remote (US timezone preferred)  
**Experience:** 7+ years  
**Reports to:** Engineering Manager

**About the Role:**
We're looking for a senior backend engineer to lead our code generation initiative. You'll design and implement the template engine that powers Aurora's code generation capabilities, enabling thousands of developers to create production-ready backends in seconds.

**Responsibilities:**
- Design and implement Aurora's template engine
- Define code generation architecture and patterns
- Mentor backend engineers on generator development
- Ensure generated code quality and best practices
- Collaborate with language specialists
- Performance optimization

**Requirements:**
- 7+ years backend development experience
- Expert in Go
- Experience with template engines or code generation
- Strong architectural skills
- API design expertise
- Production systems experience
- Excellent communication

**Nice to Have:**
- Experience with compiler design or AST manipulation
- Multiple programming language expertise
- Open source contributions
- Technical writing or speaking

**Compensation:**
- $170,000 - $220,000 base salary
- Equity
- Benefits (health, dental, vision, 401k)
- $2,000 annual learning budget

---

### Frontend Engineer

**Location:** Remote  
**Experience:** 3-6 years  
**Reports to:** Frontend Lead

**About the Role:**
Join Aurora's frontend team to build the developer dashboard that thousands will use to generate and manage their backend infrastructure. You'll work on a React/TypeScript application with a focus on developer experience.

**Responsibilities:**
- Build UI components for code generation flow
- Implement download and project management features
- Write tests for frontend code
- Collaborate with backend team on API integration
- Ensure accessibility and performance

**Requirements:**
- 3+ years frontend development
- Strong React and TypeScript skills
- CSS and responsive design
- Testing (Jest, Testing Library)
- Git workflows
- API integration experience

**Nice to Have:**
- Developer tools experience
- Design system experience
- Performance optimization
- Accessibility expertise

**Compensation:**
- $130,000 - $170,000 base salary
- Equity
- Benefits
- Remote work setup budget

---

**For full job descriptions and to apply, visit:** `careers.aurorabaas.com`

---

**Document Status:** Living Document  
**Last Updated:** May 11, 2026  
**Next Review:** Quarterly or as team grows  
**Maintained By:** Engineering Manager
