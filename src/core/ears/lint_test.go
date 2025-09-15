package ears

import "testing"

func TestParse_Ubiquitous(t *testing.T) {
	res, err := ParseRequirement("The system shall record events")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res.Shape != ShapeUbiquitous {
		t.Fatalf("expected shape ubiquitous, got %v", res.Shape)
	}
}

func TestParse_Event(t *testing.T) {
	res, err := ParseRequirement("When button is pressed, the controller shall start")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res.Shape != ShapeEvent {
		t.Fatalf("shape=%v", res.Shape)
	}
}

func TestParse_State(t *testing.T) {
	res, err := ParseRequirement("While in armed mode, the system shall alert")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res.Shape != ShapeState {
		t.Fatalf("shape=%v", res.Shape)
	}
}

func TestParse_Complex(t *testing.T) {
	res, err := ParseRequirement("While battery is low, when charger is connected, the device shall charge")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res.Shape != ShapeComplex {
		t.Fatalf("shape=%v", res.Shape)
	}
}

func TestParse_Unwanted(t *testing.T) {
	res, err := ParseRequirement("If overheating, then the system shall shutdown")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res.Shape != ShapeUnwanted {
		t.Fatalf("shape=%v", res.Shape)
	}
}

func TestParse_Invalid(t *testing.T) {
	_, err := ParseRequirement("Because of X the system might respond")
	if err == nil {
		t.Fatalf("expected error for invalid form")
	}
}

func TestExtracts_Content(t *testing.T) {
	r, err := ParseRequirement("When user clicks, the payment service shall create invoice")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if r.Trigger == "" || r.System == "" || r.Response == "" {
		t.Fatalf("expected non-empty fields: %+v", r)
	}
}

func TestMultiplePreconditions(t *testing.T) {
	r, err := ParseRequirement("While battery is low and door is open, the system shall beep")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(r.Preconditions) < 2 {
		t.Fatalf("expected >=2 preconditions, got %d", len(r.Preconditions))
	}
}

func TestLint_Batch(t *testing.T) {
	lines := []string{
		"The gateway shall route",
		"Because invalid",
		"If anomaly, then the system shall alert",
	}
	issues := Lint(lines)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Line != 2 {
		t.Fatalf("expected issue at line 2, got %d", issues[0].Line)
	}
}
