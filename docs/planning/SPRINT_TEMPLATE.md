# Aurora Sprint Planning Template

**Purpose:** Standardize sprint planning, execution, and retrospectives across Aurora teams  
**Sprint Duration:** 2 weeks  
**Last Updated:** May 11, 2026

---

## Table of Contents

- [Sprint Overview](#sprint-overview)
- [Sprint Planning Template](#sprint-planning-template)
- [Daily Standup Template](#daily-standup-template)
- [Sprint Review Template](#sprint-review-template)
- [Sprint Retrospective Template](#sprint-retrospective-template)
- [Story Point Estimation Guide](#story-point-estimation-guide)
- [Sprint Checklist](#sprint-checklist)

---

## Sprint Overview

### Sprint Cadence

```
Week 1:
  Monday    - Sprint Planning (1 hour)
  Tue-Fri   - Development + Daily Standups
  Wednesday - Mid-sprint Sync (optional, 30 min)
  
Week 2:
  Mon-Thu   - Development + Daily Standups
  Thursday  - Code freeze for demo prep
  Friday    - Sprint Review (30 min) + Retrospective (1 hour)
  Friday PM - Next sprint prep
```

### Key Principles

1. **Sustainable Pace** - No overtime, no heroics
2. **Team Commitment** - Sprint goal over individual tasks
3. **Flexibility** - Adapt to learnings and blockers
4. **Quality First** - Done means tested and documented
5. **Continuous Improvement** - Retrospectives drive change

---

## Sprint Planning Template

### Meeting Details

**Duration:** 1 hour  
**When:** Monday 10:00 AM (Week 1)  
**Who:** Full team + Product Manager  
**Location:** Zoom or conference room

### Agenda

**Part 1: Sprint Goal (15 min)**
- Review roadmap and priorities
- Define sprint goal (1-2 sentences)
- Identify any dependencies or blockers

**Part 2: Capacity Planning (10 min)**
- Calculate team capacity
- Account for PTO, holidays, etc.
- Set realistic commitment

**Part 3: Backlog Refinement (25 min)**
- Review top priority stories
- Break down large stories
- Estimate story points
- Identify unknowns or risks

**Part 4: Task Assignment (10 min)**
- Team members volunteer for stories
- Balance workload
- Identify collaboration needs

---

### Sprint Planning Document Template

```markdown
# Sprint [Number]: [Name]

**Dates:** [Start Date] - [End Date]  
**Sprint Goal:** [1-2 sentence goal]

## Team Capacity

| Team Member | Days Available | Capacity (points) | Notes |
|-------------|----------------|-------------------|-------|
| Alice       | 10             | 20                | -     |
| Bob         | 8              | 16                | PTO Fri |
| Carol       | 10             | 20                | -     |
| Dave        | 10             | 20                | -     |
| **Total**   | **38**         | **76**            | -     |

## Sprint Backlog

### P0 - Must Have (for Sprint Goal)

| Story ID | Title | Assignee | Points | Status |
|----------|-------|----------|--------|--------|
| ABC-123 | Implement TypeScript generator | Alice | 8 | Todo |
| ABC-124 | Create Express.js templates | Bob | 5 | Todo |
| ABC-125 | Add Prisma schema generation | Carol | 5 | Todo |

**P0 Total:** 18 points

### P1 - Should Have

| Story ID | Title | Assignee | Points | Status |
|----------|-------|----------|--------|--------|
| ABC-126 | Implement input validation | Dave | 3 | Todo |
| ABC-127 | Add error handling middleware | Alice | 3 | Todo |

**P1 Total:** 6 points

### P2 - Nice to Have

| Story ID | Title | Assignee | Points | Status |
|----------|-------|----------|--------|--------|
| ABC-128 | Improve logging | Bob | 2 | Todo |

**P2 Total:** 2 points

**Committed Total:** 26 points (34% of capacity - conservative)

## Dependencies

- [ ] Design mockups for code gen UI (Carol waiting)
- [ ] PostgreSQL schema finalized (Dave dependency)

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Template complexity may take longer | High | Start simple, timebox research |
| Bob out Friday | Low | Front-load critical work |

## Sprint Goal Acceptance Criteria

- [ ] TypeScript generator produces valid code
- [ ] Generated code passes linter
- [ ] End-to-end test with sample schema passes
- [ ] Documentation updated
```

---

## Daily Standup Template

### Meeting Details

**Duration:** 15 minutes (strict)  
**When:** 9:00 AM daily  
**Who:** Full team  
**Format:** Round-robin

### Structure

Each person answers:
1. **What I did yesterday**
2. **What I'll do today**
3. **Blockers** (if any)

### Example Updates

**Alice:**
> **Yesterday:** Completed TypeScript generator interface, started Express templates  
> **Today:** Finish Express templates, begin testing  
> **Blockers:** None

**Bob:**
> **Yesterday:** Researched Prisma schema generation approach  
> **Today:** Implement Prisma schema generator  
> **Blockers:** Waiting on schema format decision (need 30 min sync with Carol)

**Carol:**
> **Yesterday:** Reviewed PRs, fixed bug in schema validation  
> **Today:** Design UI mockups for code gen, sync with Bob  
> **Blockers:** None

**Dave:**
> **Yesterday:** Out sick  
> **Today:** Catch up on PRs, start validation work  
> **Blockers:** None

### Standup Guidelines

**Do:**
- Keep it brief (2 min per person max)
- Focus on sprint goal
- Call out blockers immediately
- Take detailed discussions offline
- Use standup bot for async updates if needed

**Don't:**
- Solve problems in standup
- Give detailed technical explanations
- Report to manager (update to team)
- Skip standup without notice

### Async Standup (for distributed teams)

Post in Slack by 9 AM:

```
**Yesterday:**
- Completed X
- Started Y

**Today:**
- Will finish Y
- Will start Z

**Blockers:**
- None / Need help with [specific thing]
```

---

## Sprint Review Template

### Meeting Details

**Duration:** 30 minutes  
**When:** Friday 3:00 PM (Week 2)  
**Who:** Team + stakeholders  
**Format:** Live demo

### Agenda

**Part 1: Sprint Goal Review (5 min)**
- Did we meet the sprint goal?
- What did we commit to?
- What did we deliver?

**Part 2: Demos (20 min)**
- Live demos of completed work
- Each feature owner presents
- Show, don't tell

**Part 3: Metrics & Feedback (5 min)**
- Velocity
- Completed vs committed
- Stakeholder feedback

### Sprint Review Document Template

```markdown
# Sprint [Number] Review

**Date:** [Date]  
**Sprint Goal:** [Goal]  
**Goal Met:** ✅ Yes / ❌ No / ⚠️ Partially

## Completed Stories

| Story ID | Title | Points | Demo Link |
|----------|-------|--------|-----------|
| ABC-123 | TypeScript generator | 8 | [Demo video] |
| ABC-124 | Express.js templates | 5 | [Demo video] |
| ABC-125 | Prisma schema generation | 5 | [Demo video] |
| ABC-126 | Input validation | 3 | [Demo video] |

**Completed:** 21 points  
**Committed:** 26 points  
**Completion Rate:** 81%

## Carried Over

| Story ID | Title | Points | Reason |
|----------|-------|--------|--------|
| ABC-127 | Error handling middleware | 3 | Blocked by design decision |
| ABC-128 | Improve logging | 2 | Deprioritized |

## Key Achievements

- ✅ TypeScript generator fully working
- ✅ Generated code passes all linters
- ✅ End-to-end test working
- ✅ Documentation complete

## Demos

### TypeScript Generator (Alice)
- Showed generation from YAML schema
- Generated Express.js API with CRUD
- Compiled without errors
- Ran integration tests live

**Feedback:**
- "This is amazing!" - Product Manager
- "Can we add GraphQL support?" - Stakeholder (noted for backlog)

### Prisma Integration (Carol)
- Generated Prisma schema from YAML
- Showed database migrations
- Demonstrated type safety

**Feedback:**
- "Prisma types look perfect" - Tech Lead

## Stakeholder Feedback

- Very happy with progress
- Excited to see Python generator next sprint
- Request: Add more examples to docs (added to backlog)

## Metrics

- **Velocity:** 21 points (expected: ~20 points per sprint)
- **Completion Rate:** 81%
- **Bugs Found:** 3 (all fixed)
- **PRs Merged:** 12

## Next Sprint Preview

- Python FastAPI generator
- SQLAlchemy integration
- Python client SDK
```

---

## Sprint Retrospective Template

### Meeting Details

**Duration:** 1 hour  
**When:** Friday 4:00 PM (Week 2)  
**Who:** Team only (no managers or stakeholders)  
**Facilitator:** Rotates each sprint

### Format Options

**Option 1: Start/Stop/Continue**
- What should we start doing?
- What should we stop doing?
- What should we continue doing?

**Option 2: Went Well / Could Improve / Action Items**
- What went well?
- What could be improved?
- What actions will we take?

**Option 3: Sailboat**
- Wind (what helps us go faster)
- Anchor (what slows us down)
- Rocks (risks ahead)
- Island (goal)

### Retrospective Document Template

```markdown
# Sprint [Number] Retrospective

**Date:** [Date]  
**Facilitator:** [Name]  
**Participants:** [Names]

## What Went Well 🎉

- TypeScript generator came together quickly
- Great collaboration between Alice and Bob
- Pair programming on complex template logic was very helpful
- Good code review turnaround times
- Documentation written alongside code

## What Could Be Improved 🤔

- Mid-sprint discovered template approach needed rethinking (costly)
- Some PRs sat for 2 days before review
- Not enough time for refactoring technical debt
- Standup sometimes runs over 15 minutes
- Testing took longer than estimated

## What We Learned 📚

- Template design needs prototype phase
- Estimation for new technologies should be more conservative
- Pair programming accelerates complex work
- Small PRs merge faster than large ones

## Action Items 🎯

| Action | Owner | Due Date | Status |
|--------|-------|----------|--------|
| Create template design doc before implementation | Alice | Next sprint planning | 📋 Todo |
| Set up code review SLA reminder bot | Dave | This Friday | 📋 Todo |
| Schedule 10% time for tech debt each sprint | Team | Next sprint | 📋 Todo |
| Timebox standup strictly to 15 min | Facilitator | Ongoing | 📋 Todo |
| Add buffer time to estimates for unknowns | Everyone | Next sprint | 📋 Todo |

## Team Mood Check 😊

Rate 1-5 (1=bad, 5=great):

- Work-life balance: 4.5/5
- Team collaboration: 5/5
- Code quality: 4/5
- Progress toward goals: 4.5/5
- Overall satisfaction: 4.5/5

**Comments:**
- "Great sprint, felt productive"
- "Enjoyed pairing with Bob, learned a lot"
- "Would like more time for testing"

## Shout-outs 👏

- Alice for excellent template design
- Bob for patient pair programming
- Carol for thorough code reviews
- Dave for jumping in to help with testing

## Next Sprint Focus Areas

- Apply learnings about template prototyping
- Improve PR turnaround with reminders
- Allocate time for technical debt
- Maintain strong collaboration
```

---

## Story Point Estimation Guide

### Fibonacci Scale

Aurora uses Fibonacci sequence for story points: **1, 2, 3, 5, 8, 13**

### Point Guidelines

| Points | Time Estimate | Complexity | Uncertainty | Examples |
|--------|---------------|------------|-------------|----------|
| **1** | 2-4 hours | Very simple | Clear solution | Minor bug fix, small refactor |
| **2** | Half day | Simple | Well-understood | Add validation rule, update docs |
| **3** | 1 day | Moderate | Some unknowns | New API endpoint with tests |
| **5** | 2-3 days | Complex | Multiple approaches | Feature with DB changes |
| **8** | 4-5 days | Very complex | Significant unknowns | New major component |
| **13** | 1+ week | Too large | Many unknowns | **SPLIT INTO SMALLER STORIES** |

### Estimation Factors

**Complexity:**
- Simple logic vs complex algorithms
- Known patterns vs new approaches
- Number of components involved

**Uncertainty:**
- How well understood is the requirement?
- Are there unknowns in implementation?
- Need research or spikes?

**Effort:**
- Development time
- Testing time
- Documentation time
- Code review time

**Dependencies:**
- External team dependencies
- Third-party library unknowns
- Infrastructure requirements

### Estimation Process

**Planning Poker:**
1. Read story aloud
2. Clarify questions
3. Everyone votes simultaneously (show cards)
4. Discuss outliers
5. Re-vote if needed
6. Consensus or average

**Reference Stories:**
- Keep examples of each point value
- "This is similar to X which was 5 points"
- Build team calibration over time

### Story Size Guidelines

**Good Story:**
- Completable within a sprint
- Has clear acceptance criteria
- Deliverable independently
- Testable
- Demonstrable

**Story Too Large:**
- Takes more than 5 days
- Spans multiple components
- Unclear how to complete
- **Action:** Split into smaller stories

**Story Too Small:**
- Takes less than 2 hours
- Just a task, not a user story
- **Action:** Combine with related stories or make it a subtask

### Examples from Aurora

**1 Point:**
```
As a developer
I want schema validation to check for reserved words
So that generated code doesn't have naming conflicts

Acceptance Criteria:
- Add list of reserved words per language
- Validation rejects reserved words
- Error message is helpful
```

**3 Points:**
```
As a user
I want to download generated code as a ZIP file
So that I can use it in my project

Acceptance Criteria:
- API endpoint returns ZIP file
- ZIP contains all generated files
- README.md included
- File structure is correct
- Works for both TypeScript and Python
```

**5 Points:**
```
As a developer
I want TypeScript Express.js code generation
So that I can get a working REST API

Acceptance Criteria:
- Generate server.ts with Express setup
- Generate routes for each table
- Generate controllers with CRUD operations
- Generate Prisma schema
- Include package.json with dependencies
- Generated code compiles without errors
- Generated code passes linter
```

**8 Points:**
```
As a user
I want authentication and authorization in generated code
So that my API is secure

Acceptance Criteria:
- JWT token generation and validation
- User registration and login endpoints
- Password hashing
- Auth middleware
- Protected routes
- Client SDK auth methods
- Tests for auth flow
- Documentation

Note: This is probably 13 points - should split into:
- JWT implementation (5 points)
- User auth endpoints (3 points)
- Auth middleware and protection (3 points)
```

---

## Sprint Checklist

### Sprint Planning Checklist

**Before the Meeting:**
- [ ] Product Manager prepares prioritized backlog
- [ ] Stories have clear acceptance criteria
- [ ] Dependencies identified
- [ ] Design work complete for top stories

**During the Meeting:**
- [ ] Sprint goal defined
- [ ] Team capacity calculated
- [ ] Stories estimated and assigned
- [ ] Dependencies noted
- [ ] Risks identified
- [ ] Everyone commits to sprint goal

**After the Meeting:**
- [ ] Sprint board updated (Jira/Linear)
- [ ] Sprint goal posted to Slack
- [ ] Calendar invites sent
- [ ] Dependencies communicated

### Daily Standup Checklist

**Before Standup:**
- [ ] Review yesterday's work
- [ ] Plan today's work
- [ ] Identify blockers

**During Standup:**
- [ ] Update in < 2 minutes
- [ ] Call out blockers
- [ ] Offer/request help

**After Standup:**
- [ ] Schedule any needed sync-ups
- [ ] Unblock teammates if possible
- [ ] Update Jira board

### Sprint Review Checklist

**Before the Meeting:**
- [ ] Demo environment ready
- [ ] Demo script prepared
- [ ] Screenshots/videos ready (backup)
- [ ] Stakeholders invited
- [ ] Metrics gathered

**During the Meeting:**
- [ ] Show working software (not slides!)
- [ ] Explain what was built and why
- [ ] Collect feedback
- [ ] Note requests for backlog

**After the Meeting:**
- [ ] Add feedback to backlog
- [ ] Share demo recording
- [ ] Update roadmap if needed

### Sprint Retrospective Checklist

**Before the Meeting:**
- [ ] Book private room/space
- [ ] Prepare retro board (Miro/Mural)
- [ ] Review last retro action items
- [ ] Choose facilitation format

**During the Meeting:**
- [ ] Review previous action items
- [ ] Gather feedback (what went well, etc.)
- [ ] Vote on top issues
- [ ] Define action items
- [ ] Assign owners and dates

**After the Meeting:**
- [ ] Document retro notes
- [ ] Add action items to next sprint
- [ ] Share summary (if appropriate)
- [ ] Follow up on action items

### Definition of Done Checklist

A story is "Done" when:

**Code:**
- [ ] Code written and working
- [ ] Code reviewed and approved
- [ ] No compiler warnings
- [ ] No linter errors
- [ ] Follows code style guide

**Testing:**
- [ ] Unit tests written and passing
- [ ] Integration tests written and passing
- [ ] Manual testing completed
- [ ] Edge cases covered
- [ ] No new bugs introduced

**Documentation:**
- [ ] Code comments added
- [ ] API documentation updated
- [ ] User documentation updated (if user-facing)
- [ ] README updated (if needed)

**Quality:**
- [ ] Performance acceptable
- [ ] Security considerations addressed
- [ ] Accessibility checked (frontend)
- [ ] Error handling present

**Integration:**
- [ ] Merged to main/develop branch
- [ ] CI/CD pipeline passing
- [ ] Deployed to staging
- [ ] Verified in staging environment

**Communication:**
- [ ] Demo prepared
- [ ] Stakeholders notified (if needed)
- [ ] Jira ticket updated
- [ ] Sprint board updated

---

## Velocity Tracking

### Calculating Velocity

**Velocity** = Average story points completed per sprint

**Example:**
```
Sprint 1: 18 points
Sprint 2: 22 points
Sprint 3: 20 points

Average velocity = (18 + 22 + 20) / 3 = 20 points/sprint
```

### Using Velocity

**Capacity Planning:**
- Use velocity to predict sprint capacity
- Don't over-commit (aim for 80-90% of velocity)
- Account for unknowns and risks

**Roadmap Planning:**
- Estimate project: 120 points
- Velocity: 20 points/sprint
- Expected duration: 6 sprints (12 weeks)

**Improving Velocity:**
- Remove blockers
- Improve estimation accuracy
- Reduce technical debt
- Better collaboration
- **Don't** pressure team to increase velocity artificially

### Velocity Anti-Patterns

**Don't:**
- Use velocity to compare teams
- Use velocity as performance metric
- Pressure team to increase velocity
- Game the system (inflate estimates)
- Focus on velocity over value

**Do:**
- Use velocity for planning
- Track trends over time
- Investigate sudden changes
- Focus on delivering value

---

## Sprint Anti-Patterns to Avoid

### 1. Overcommitment
**Problem:** Team commits to too much work  
**Symptoms:** Incomplete stories, rushed work, low quality  
**Solution:** Commit to 70-80% of capacity, leave buffer

### 2. Scope Creep
**Problem:** Stories added mid-sprint  
**Symptoms:** Sprint goal at risk, team distracted  
**Solution:** Protect sprint, add to next sprint unless critical

### 3. Hero Culture
**Problem:** One person saving the sprint with overtime  
**Symptoms:** Burnout, single point of failure  
**Solution:** Sustainable pace, team collaboration

### 4. Zombie Stories
**Problem:** Stories started but never finished  
**Symptoms:** Many "In Progress", nothing "Done"  
**Solution:** Limit WIP, focus on completion over starting

### 5. Meeting Fatigue
**Problem:** Too many or too long meetings  
**Symptoms:** No time for actual work  
**Solution:** Strict timeboxes, skip if not valuable

### 6. No Retrospective Action
**Problem:** Retro action items never implemented  
**Symptoms:** Same issues every sprint  
**Solution:** Make action items stories, track them

### 7. Technical Debt Neglect
**Problem:** Never time for refactoring or cleanup  
**Symptoms:** Velocity decreases over time  
**Solution:** Allocate 10-20% capacity for tech debt

### 8. Solo Work
**Problem:** Everyone works in isolation  
**Symptoms:** Knowledge silos, delayed reviews  
**Solution:** Pair programming, mob programming, collaboration

---

## Sprint Metrics Dashboard

Track these metrics per sprint:

### Velocity Metrics
- **Planned Points:** Points committed at sprint start
- **Completed Points:** Points actually completed
- **Completion Rate:** Completed / Planned (target: 80-100%)
- **Velocity:** Average completed points (track trend)

### Quality Metrics
- **Bugs Created:** New bugs introduced
- **Bugs Fixed:** Bugs resolved
- **Test Coverage:** % of code covered by tests
- **Code Review Time:** Average time for review

### Process Metrics
- **Cycle Time:** Time from start to done
- **WIP:** Number of stories in progress
- **Blocked Days:** Days stories were blocked
- **Carryover:** Stories moved to next sprint

### Team Health Metrics
- **Team Mood:** Average satisfaction (1-5)
- **Collaboration:** Pair programming sessions
- **Learning:** Tech talks, knowledge sharing
- **Balance:** PTO days, overtime hours

---

## Appendix: Templates for Copy/Paste

### Sprint Goal Template
```
Sprint [N] Goal:
[Deliver X capability that enables Y value for users]

Success Criteria:
- [ ] [Measurable outcome 1]
- [ ] [Measurable outcome 2]
- [ ] [Measurable outcome 3]
```

### User Story Template
```
As a [user type]
I want [capability]
So that [benefit/value]

Acceptance Criteria:
- [ ] [Specific, testable criterion 1]
- [ ] [Specific, testable criterion 2]
- [ ] [Specific, testable criterion 3]

Technical Notes:
[Any implementation details, constraints, or considerations]

Definition of Done:
- [ ] Code complete and reviewed
- [ ] Tests passing (unit + integration)
- [ ] Documentation updated
- [ ] Deployed to staging
```

### Bug Report Template
```
**Bug:** [Short description]

**Severity:** P0 (Critical) / P1 (High) / P2 (Medium) / P3 (Low)

**Environment:**
- OS: [e.g., macOS 12]
- Browser: [e.g., Chrome 104]
- Version: [e.g., v0.1.0]

**Steps to Reproduce:**
1. [Step 1]
2. [Step 2]
3. [Step 3]

**Expected Result:**
[What should happen]

**Actual Result:**
[What actually happened]

**Screenshots/Logs:**
[Attach if applicable]

**Impact:**
[Who is affected? How many users?]
```

### Demo Script Template
```
**Feature:** [Name]
**Demo by:** [Name]
**Duration:** [5 min]

**Setup:**
- Open [URL/App]
- Logged in as [user type]
- [Any other setup]

**Demo Flow:**
1. Show problem/context
   - [What problem does this solve?]
2. Demonstrate feature
   - [Step 1: Do X, see Y]
   - [Step 2: Do A, see B]
3. Highlight key points
   - [Performance, UX, etc.]
4. Show edge cases
   - [Error handling, validation]

**Call to Action:**
- [What do we want from stakeholders?]
- [Questions to ask?]
```

---

## Quick Reference

### Sprint at a Glance

```
Sprint Duration: 2 weeks

Week 1:
  Mon: Planning (1h)
  Tue-Fri: Build

Week 2:
  Mon-Thu: Build
  Fri: Review (30m) + Retro (1h)

Daily: Standup (15m @ 9am)

Point Scale: 1, 2, 3, 5, 8, 13
Target Velocity: ~20 points/sprint

Commit to: 70-80% of capacity
```

### Common Commands

**Start sprint:**
```bash
# Create sprint in Jira
# Post sprint goal to #aurora-dev
# Update sprint board
```

**Daily update:**
```
Yesterday: [X]
Today: [Y]
Blockers: [Z or None]
```

**End sprint:**
```bash
# Run sprint report
# Prepare demos
# Schedule retro
# Archive completed stories
```

---

**Document Owner:** Scrum Master / Engineering Manager  
**Last Updated:** May 11, 2026  
**Next Review:** After 3 sprints or as needed  
**Feedback:** #aurora-dev or retro meetings
