package codescan

import "testing"

// TGPlatform billing must make 3.1.1 evaluable on its own. The SDK's purchase
// call happens to contain "PurchaseProduct", which the generic pattern already
// matches — but relying on that is relying on a coincidence of naming. These
// cases pin the namespace-anchored detection so a storefront that only
// initialises or queries billing still triggers the rule.
func TestTGPlatformBillingTriggersIAPRule(t *testing.T) {
	r := ruleByID(t, "iap-no-restore")

	triggers := []string{
		`TGP.Billing.SetUpBilling(config, OnReady);`,
		`TGP.Billing.QueryInAppProduct(ids, OnList);`,
		`TGP.Billing.IAP.PurchaseProduct(productId, OnDone);`,
		`TGP.Billing.IAP.ConsumeProduct(receipt);`,
	}
	for _, line := range triggers {
		if got := r.Check(csCtx(line)); len(got) == 0 {
			t.Errorf("TGP billing should trigger iap-no-restore: %q", line)
		}
	}

	// The SDK's restore call satisfies the anti-pattern, so a storefront that
	// implements restore is not flagged.
	restore := csCtx(`var purchaseReceipt = TGP.Billing.IAP.RestorePurchase();`)
	if !r.AntiPatternMatched(restore) {
		t.Error("TGP.Billing.IAP.RestorePurchase() must count as a restore implementation")
	}
}

// DetectClaims must see TGPlatform billing too, so the IAP flow is claimed for
// projects that use it.
func TestTGPlatformBillingClaimsIAP(t *testing.T) {
	if !reIAP.MatchString(`TGP.Billing.SetUpBilling(config, OnReady);`) {
		t.Error("TGP billing should claim the IAP flow")
	}
}
