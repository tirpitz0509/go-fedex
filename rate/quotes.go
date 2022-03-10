package rate

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"strings"

	conv "github.com/cstockton/go-conv"
	"github.com/tirpitz0509/go-fedex/common"
)

const (
	apiWSDLTest = "https://wsbeta.fedex.com:443/web-services/rate"
	apiWSDLLive = "https://ws.fedex.com:443/web-services/rate"
	RATESOAP    = "http://fedex.com/ws/rate/v28"
)

type Address struct {
	StreetLines         []string `json:"streetLines,omitempty"`         //
	City                string   `json:"city,omitempty"`                //
	StateOrProvinceCode string   `json:"stateOrProvinceCode,omitempty"` //
	PostalCode          string   `json:"postalCode"`                    //
	CountryCode         string   `json:"countryCode"`                   //
	Residential         bool     `json:"residential,omitempty"`         //
}

type Package struct {
	SubPackagingType  string `json:"subPackagingType,omitempty"`
	GroupPackageCount int    `json:"groupPackageCount,omitempty"`
	ContentRecord     []struct {
		ItemNumber       string `json:"itemNumber"`
		ReceivedQuantity int    `json:"receivedQuantity"`
		Description      string `json:"description"`
		PartNumber       string `json:"partNumber"`
	} `json:"contentRecord,omitempty"`
	DeclaredValue struct {
		Amount   string `json:"amount,omitempty"`
		Currency string `json:"currency,omitempty"`
	} `json:"declaredValue,omitempty"`
	Weight struct {
		Units string `json:"units"`
		Value int    `json:"value"`
	} `json:"weight"`
	Dimensions struct {
		Length int    `json:"length,omitempty"`
		Width  int    `json:"width,omitempty"`
		Height int    `json:"height,omitempty"`
		Units  string `json:"units,omitempty"`
	} `json:"dimensions,omitempty"`
	VariableHandlingChargeDetail struct {
		RateType      string `json:"rateType,omitempty"`
		PercentValue  int    `json:"percentValue,omitempty"`
		RateLevelType string `json:"rateLevelType,omitempty"`
		FixedValue    struct {
			Amount   string `json:"amount,omitempty"`
			Currency string `json:"currency,omitempty"`
		} `json:"fixedValue,omitempty"`
		RateElementBasis string `json:"rateElementBasis,omitempty"`
	} `json:"variableHandlingChargeDetail,omitempty"`
	PackageSpecialServices struct {
		SpecialServiceTypes []string `json:"specialServiceTypes,omitempty"`
		AlcoholDetail       struct {
			AlcoholRecipientType string `json:"alcoholRecipientType,omitempty"`
			ShipperAgreementType string `json:"shipperAgreementType,omitempty"`
		} `json:"alcoholDetail,omitempty"`
		DangerousGoodsDetail struct {
			Offeror                string   `json:"offeror,omitempty"`
			Accessibility          string   `json:"accessibility,omitempty"`
			EmergencyContactNumber string   `json:"emergencyContactNumber,omitempty"`
			Options                []string `json:"options,omitempty"`
			Containers             []struct {
				Offeror              string `json:"offeror,omitempty"`
				HazardousCommodities []struct {
					Quantity struct {
						QuantityType string `json:"quantityType,omitempty"`
						Amount       int    `json:"amount,omitempty"`
						Units        string `json:"units,omitempty"`
					} `json:"quantity,omitempty"`
					InnerReceptacles []struct {
						Quantity struct {
							QuantityType string `json:"quantityType,omitempty"`
							Amount       int    `json:"amount,omitempty"`
							Units        string `json:"units,omitempty"`
						} `json:"quantity,omitempty"`
					} `json:"innerReceptacles,omitempty"`
					Options struct {
						LabelTextOption           string `json:"labelTextOption,omitempty"`
						CustomerSuppliedLabelText string `json:"customerSuppliedLabelText,omitempty"`
					} `json:"options,omitempty"`
					Description struct {
						SequenceNumber    int      `json:"sequenceNumber,omitempty"`
						ProcessingOptions []string `json:"processingOptions,omitempty"`
						SubsidiaryClasses string   `json:"subsidiaryClasses,omitempty"`
						LabelText         string   `json:"labelText,omitempty"`
						TechnicalName     string   `json:"technicalName,omitempty"`
						PackingDetails    struct {
							PackingInstructions string `json:"packingInstructions,omitempty"`
							CargoAircraftOnly   bool   `json:"cargoAircraftOnly,omitempty"`
						} `json:"packingDetails,omitempty"`
						Authorization      string `json:"authorization,omitempty"`
						ReportableQuantity bool   `json:"reportableQuantity,omitempty"`
						Percentage         int    `json:"percentage,omitempty"`
						ID                 string `json:"id,omitempty"`
						PackingGroup       string `json:"packingGroup,omitempty"`
						ProperShippingName string `json:"properShippingName,omitempty"`
						HazardClass        string `json:"hazardClass,omitempty"`
					} `json:"description,omitempty"`
				} `json:"hazardousCommodities,omitempty"`
				NumberOfContainers     int    `json:"numberOfContainers,omitempty"`
				ContainerType          string `json:"containerType,omitempty"`
				EmergencyContactNumber struct {
					AreaCode                     string `json:"areaCode,omitempty"`
					Extension                    string `json:"extension,omitempty"`
					CountryCode                  string `json:"countryCode,omitempty"`
					PersonalIdentificationNumber string `json:"personalIdentificationNumber,omitempty"`
					LocalNumber                  string `json:"localNumber,omitempty"`
				} `json:"emergencyContactNumber,omitempty"`
				Packaging struct {
					Count int    `json:"count,omitempty"`
					Units string `json:"units,omitempty"`
				} `json:"packaging,omitempty"`
				PackingType               string `json:"packingType,omitempty"`
				RadioactiveContainerClass string `json:"radioactiveContainerClass,omitempty"`
			} `json:"containers,omitempty"`
			Packaging struct {
				Count int    `json:"count,omitempty"`
				Units string `json:"units,omitempty"`
			} `json:"packaging,omitempty"`
		} `json:"dangerousGoodsDetail,omitempty"`
		PackageCODDetail struct {
			CodCollectionAmount struct {
				Amount   float64 `json:"amount,omitempty"`
				Currency string  `json:"currency,omitempty"`
			} `json:"codCollectionAmount,omitempty"`
			CodCollectionType string `json:"codCollectionType,omitempty"`
		} `json:"packageCODDetail,omitempty"`
		PieceCountVerificationBoxCount int `json:"pieceCountVerificationBoxCount,omitempty"`
		BatteryDetails                 []struct {
			Material          string `json:"material,omitempty"`
			RegulatorySubType string `json:"regulatorySubType,omitempty"`
			Packing           string `json:"packing,omitempty"`
		} `json:"batteryDetails,omitempty"`
		DryIceWeight struct {
			Units string `json:"units,omitempty"`
			Value int    `json:"value,omitempty"`
		} `json:"dryIceWeight,omitempty"`
	} `json:"packageSpecialServices,omitempty"`
}

type RateRequest struct {
	AccountNumber struct {
		Value string `json:"value"` //
	} `json:"accountNumber"` //
	RateRequestControlParameters struct {
		ReturnTransitTimes          bool   `json:"returnTransitTimes,omitempty"`          //
		ServicesNeededOnRateFailure bool   `json:"servicesNeededOnRateFailure,omitempty"` //
		VariableOptions             string `json:"variableOptions,omitempty"`             //
		RateSortOrder               string `json:"rateSortOrder,omitempty"`               //
	} `json:"rateRequestControlParameters,omitempty"` //
	RequestedShipment struct {
		Shipper struct {
			Address Address `json:"address"` //
		} `json:"shipper,omitempty"` //
		Recipient struct {
			Address Address `json:"address"` //
		} `json:"recipient,omitempty"` //
		ServiceType             string `json:"serviceType"` //
		EmailNotificationDetail struct {
			Recipients []struct {
				EmailAddress          string   `json:"emailAddress,omitempty"`
				NotificationEventType []string `json:"notificationEventType,omitempty"`
				SmsDetail             struct {
					PhoneNumber            string `json:"phoneNumber,omitempty"`
					PhoneNumberCountryCode string `json:"phoneNumberCountryCode,omitempty"`
				} `json:"smsDetail,omitempty"`
				NotificationFormatType         string `json:"notificationFormatType,omitempty"`
				EmailNotificationRecipientType string `json:"emailNotificationRecipientType,omitempty"`
				NotificationType               string `json:"notificationType,omitempty"`
				Locale                         string `json:"locale,omitempty"`
			} `json:"recipients,omitempty"`
			PersonalMessage  string `json:"personalMessage,omitempty"`
			PrintedReference struct {
				PrintedReferenceType string `json:"printedReferenceType,omitempty"`
				Value                string `json:"value,omitempty"`
			} `json:"PrintedReference,omitempty"`
		} `json:"emailNotificationDetail,omitempty"` //
		PreferredCurrency         string    `json:"preferredCurrency,omitempty"`         //
		RateRequestType           []string  `json:"rateRequestType,omitempty"`           //
		ShipDateStamp             string    `json:"shipDateStamp,omitempty"`             //
		PickupType                string    `json:"pickupType,omitempty"`                //
		RequestedPackageLineItems []Package `json:"requestedPackageLineItems,omitempty"` //
		DocumentShipment          bool      `json:"documentShipment,omitempty"`          //
		PickupDetail              struct {
			CompanyCloseTime string `json:"companyCloseTime,omitempty"`
			PickupOrigin     struct {
				AccountNumber struct {
					Value int `json:"value,omitempty"`
				} `json:"accountNumber,omitempty"`
				Address struct {
					AddressVerificationID string   `json:"addressVerificationId,omitempty"`
					CountryCode           string   `json:"countryCode,omitempty"`
					StreetLines           []string `json:"streetLines,omitempty"`
				} `json:"address,omitempty"`
				Contact struct {
					CompanyName string `json:"companyName,omitempty"`
					FaxNumber   string `json:"faxNumber,omitempty"`
					PersonName  string `json:"personName,omitempty"`
					PhoneNumber string `json:"phoneNumber,omitempty"`
				} `json:"contact,omitempty"`
			} `json:"pickupOrigin,omitempty"`
			GeographicalPostalCode  string `json:"geographicalPostalCode,omitempty"`
			RequestType             string `json:"requestType,omitempty"`
			BuildingPartDescription string `json:"buildingPartDescription,omitempty"`
			CourierInstructions     string `json:"courierInstructions,omitempty"`
			BuildingPart            string `json:"buildingPart,omitempty"`
			LatestPickupDateTime    string `json:"latestPickupDateTime,omitempty"`
			PackageLocation         string `json:"packageLocation,omitempty"`
			ReadyPickupDateTime     string `json:"readyPickupDateTime,omitempty"`
			EarlyPickup             bool   `json:"earlyPickup,omitempty"`
		} `json:"pickupDetail,omitempty"` //
		VariableHandlingChargeDetail struct {
			RateType      string `json:"rateType,omitempty"`
			PercentValue  int    `json:"percentValue,omitempty"`
			RateLevelType string `json:"rateLevelType,omitempty"`
			FixedValue    struct {
				Amount   string `json:"amount,omitempty"`
				Currency string `json:"currency,omitempty"`
			} `json:"fixedValue,omitempty"`
			RateElementBasis string `json:"rateElementBasis,omitempty"`
		} `json:"variableHandlingChargeDetail,omitempty"` //
		PackagingType           string  `json:"packagingType,omitempty"`     //
		TotalPackageCount       int     `json:"totalPackageCount,omitempty"` //
		TotalWeight             float64 `json:"totalWeight,omitempty"`       //
		ShipmentSpecialServices struct {
			ReturnShipmentDetail struct {
				ReturnType string `json:"returnType,omitempty"`
			} `json:"returnShipmentDetail,omitempty"`
			DeliveryOnInvoiceAcceptanceDetail struct {
				Recipient struct {
					AccountNumber struct {
						Value int `json:"value,omitempty"`
					} `json:"accountNumber,omitempty"`
					Address struct {
						StreetLines []string `json:"streetLines,omitempty"`
						CountryCode string   `json:"countryCode,omitempty"`
					} `json:"address,omitempty"`
					Contact struct {
						CompanyName string `json:"companyName,omitempty"`
						FaxNumber   string `json:"faxNumber,omitempty"`
						PersonName  string `json:"personName,omitempty"`
						PhoneNumber string `json:"phoneNumber,omitempty"`
					} `json:"contact,omitempty"`
				} `json:"recipient,omitempty"`
			} `json:"deliveryOnInvoiceAcceptanceDetail,omitempty"`
			InternationalTrafficInArmsRegulationsDetail struct {
				LicenseOrExemptionNumber string `json:"licenseOrExemptionNumber,omitempty"`
			} `json:"internationalTrafficInArmsRegulationsDetail,omitempty"`
			PendingShipmentDetail struct {
				PendingShipmentType string `json:"pendingShipmentType,omitempty"`
				ProcessingOptions   struct {
					Options []string `json:"options,omitempty"`
				} `json:"processingOptions,omitempty"`
				RecommendedDocumentSpecification struct {
					Types []string `json:"types,omitempty"`
				} `json:"recommendedDocumentSpecification,omitempty"`
				EmailLabelDetail struct {
					Recipients []struct {
						EmailAddress     string `json:"emailAddress,omitempty"`
						OptionsRequested struct {
							Options []string `json:"options,omitempty"`
						} `json:"optionsRequested,omitempty"`
						Role   string `json:"role,omitempty"`
						Locale struct {
							Country  string `json:"country,omitempty"`
							Language string `json:"language,omitempty"`
						} `json:"locale,omitempty"`
					} `json:"recipients,omitempty"`
					Message string `json:"message,omitempty"`
				} `json:"emailLabelDetail,omitempty"`
				DocumentReferences []struct {
					DocumentType      string `json:"documentType,omitempty"`
					CustomerReference string `json:"customerReference,omitempty"`
					Description       string `json:"description,omitempty"`
					DocumentID        string `json:"documentId,omitempty"`
				} `json:"documentReferences,omitempty"`
				ExpirationTimeStamp  string `json:"expirationTimeStamp,omitempty"`
				ShipmentDryIceDetail struct {
					TotalWeight struct {
						Units string `json:"units,omitempty"`
						Value int    `json:"value,omitempty"`
					} `json:"totalWeight,omitempty"`
					PackageCount int `json:"packageCount,omitempty"`
				} `json:"shipmentDryIceDetail,omitempty"`
			} `json:"pendingShipmentDetail,omitempty"`
			HoldAtLocationDetail struct {
				LocationID                string `json:"locationId,omitempty"`
				LocationContactAndAddress struct {
					Address struct {
						StreetLines         []string `json:"streetLines,omitempty"`
						City                string   `json:"city,omitempty"`
						StateOrProvinceCode string   `json:"stateOrProvinceCode,omitempty"`
						PostalCode          string   `json:"postalCode,omitempty"`
						CountryCode         string   `json:"countryCode,omitempty"`
						Residential         bool     `json:"residential,omitempty"`
					} `json:"address,omitempty"`
					Contact struct {
						PersonName       string `json:"personName,omitempty"`
						EmailAddress     string `json:"emailAddress,omitempty"`
						ParsedPersonName struct {
							FirstName  string `json:"firstName,omitempty"`
							LastName   string `json:"lastName,omitempty"`
							MiddleName string `json:"middleName,omitempty"`
							Suffix     string `json:"suffix,omitempty"`
						} `json:"parsedPersonName,omitempty"`
						PhoneNumber    string `json:"phoneNumber,omitempty"`
						PhoneExtension string `json:"phoneExtension,omitempty"`
						CompanyName    string `json:"companyName,omitempty"`
						FaxNumber      string `json:"faxNumber,omitempty"`
					} `json:"contact,omitempty"`
				} `json:"locationContactAndAddress,omitempty"`
				LocationType string `json:"locationType,omitempty"`
			} `json:"holdAtLocationDetail,omitempty"`
			ShipmentCODDetail struct {
				AddTransportationChargesDetail struct {
					RateType        string `json:"rateType,omitempty"`
					RateLevelType   string `json:"rateLevelType,omitempty"`
					ChargeLevelType string `json:"chargeLevelType,omitempty"`
					ChargeType      string `json:"chargeType,omitempty"`
				} `json:"addTransportationChargesDetail,omitempty"`
				CodRecipient struct {
					AccountNumber struct {
						Value int `json:"value,omitempty"`
					} `json:"accountNumber,omitempty"`
				} `json:"codRecipient"`
				RemitToName                           string `json:"remitToName,omitempty"`
				CodCollectionType                     string `json:"codCollectionType,omitempty"`
				FinancialInstitutionContactAndAddress struct {
					Address struct {
						StreetLines         []string `json:"streetLines,omitempty"`
						City                string   `json:"city,omitempty"`
						StateOrProvinceCode string   `json:"stateOrProvinceCode,omitempty"`
						PostalCode          string   `json:"postalCode,omitempty"`
						CountryCode         string   `json:"countryCode,omitempty"`
						Residential         bool     `json:"residential,omitempty"`
					} `json:"address,omitempty"`
					Contact struct {
						PersonName       string `json:"personName,omitempty"`
						EmailAddress     string `json:"emailAddress,omitempty"`
						ParsedPersonName struct {
							FirstName  string `json:"firstName,omitempty"`
							LastName   string `json:"lastName,omitempty"`
							MiddleName string `json:"middleName,omitempty"`
							Suffix     string `json:"suffix,omitempty"`
						} `json:"parsedPersonName,omitempty"`
						PhoneNumber    string `json:"phoneNumber,omitempty"`
						PhoneExtension string `json:"phoneExtension,omitempty"`
						CompanyName    string `json:"companyName,omitempty"`
						FaxNumber      string `json:"faxNumber,omitempty"`
					} `json:"contact,omitempty"`
				} `json:"financialInstitutionContactAndAddress,omitempty"`
				ReturnReferenceIndicatorType string `json:"returnReferenceIndicatorType,omitempty"`
			} `json:"shipmentCODDetail,omitempty"`
			ShipmentDryIceDetail struct {
				TotalWeight struct {
					Units string `json:"units,omitempty"`
					Value int    `json:"value,omitempty"`
				} `json:"totalWeight,omitempty"`
				PackageCount int `json:"packageCount,omitempty"`
			} `json:"shipmentDryIceDetail,omitempty"`
			InternationalControlledExportDetail struct {
				Type string `json:"type,omitempty"`
			} `json:"internationalControlledExportDetail,omitempty"`
			HomeDeliveryPremiumDetail struct {
				PhoneNumber struct {
					AreaCode                     string `json:"areaCode,omitempty"`
					Extension                    string `json:"extension,omitempty"`
					CountryCode                  string `json:"countryCode,omitempty"`
					PersonalIdentificationNumber string `json:"personalIdentificationNumber,omitempty"`
					LocalNumber                  string `json:"localNumber,omitempty"`
				} `json:"phoneNumber,omitempty"`
				ShipTimestamp           string `json:"shipTimestamp,omitempty"`
				HomedeliveryPremiumType string `json:"homedeliveryPremiumType,omitempty"`
			} `json:"homeDeliveryPremiumDetail,omitempty"`
			SpecialServiceTypes []string `json:"specialServiceTypes,omitempty"`
		} `json:"shipmentSpecialServices,omitempty"` //
		CustomsClearanceDetail struct {
			CommercialInvoice struct {
				ShipmentPurpose string `json:"shipmentPurpose,omitempty"`
			} `json:"commercialInvoice,omitempty"`
			FreightOnValue string `json:"freightOnValue,omitempty"`
			DutiesPayment  struct {
				Payor struct {
					ResponsibleParty struct {
						Address struct {
							StreetLines         []string `json:"streetLines,omitempty"`
							City                string   `json:"city,omitempty"`
							StateOrProvinceCode string   `json:"stateOrProvinceCode,omitempty"`
							PostalCode          string   `json:"postalCode,omitempty"`
							CountryCode         string   `json:"countryCode,omitempty"`
							Residential         bool     `json:"residentia,omitempty"`
						} `json:"address,omitempty"`
						Contact struct {
							PersonName       string `json:"personName,omitempty"`
							EmailAddress     string `json:"emailAddress,omitempty"`
							ParsedPersonName struct {
								FirstName  string `json:"firstName,omitempty"`
								LastName   string `json:"lastName,omitempty"`
								MiddleName string `json:"middleName,omitempty"`
								Suffix     string `json:"suffix,omitempty"`
							} `json:"parsedPersonName,omitempty"`
							PhoneNumber    string `json:"phoneNumber,omitempty"`
							PhoneExtension string `json:"phoneExtension,omitempty"`
							CompanyName    string `json:"companyName,omitempty"`
							FaxNumber      string `json:"faxNumber,omitempty"`
						} `json:"contact,omitempty"`
						AccountNumber struct {
							Value string `json:"value,omitempty"`
						} `json:"accountNumber,omitempty"`
					} `json:"responsibleParty,omitempty"`
				} `json:"payor,omitempty"`
				PaymentType string `json:"paymentType,omitempty"`
			} `json:"dutiesPayment,omitempty"`
			Commodities []struct {
				Description string `json:"description,omitempty"`
				Weight      struct {
					Units string `json:"units,omitempty"`
					Value int    `json:"value,omitempty"`
				} `json:"weight,omitempty"`
				Quantity     int `json:"quantity,omitempty"`
				CustomsValue struct {
					Amount   string `json:"amount,omitempty"`
					Currency string `json:"currency,omitempty"`
				} `json:"customsValue,omitempty"`
				UnitPrice struct {
					Amount   string `json:"amount,omitempty"`
					Currency string `json:"currency,omitempty"`
				} `json:"unitPrice,omitempty"`
				NumberOfPieces       int    `json:"numberOfPieces,omitempty"`
				CountryOfManufacture string `json:"countryOfManufacture,omitempty"`
				QuantityUnits        string `json:"quantityUnits,omitempty"`
				Name                 string `json:"name,omitempty"`
				HarmonizedCode       string `json:"harmonizedCode,omitempty"`
				PartNumber           string `json:"partNumber,omitempty"`
			} `json:"commodities,omitempty"`
		} `json:"customsClearanceDetail,omitempty"` //
		GroupShipment     bool `json:"groupShipment,omitempty"` //
		ServiceTypeDetail struct {
			CarrierCode     string `json:"carrierCode,omitempty"`
			Description     string `json:"description,omitempty"`
			ServiceName     string `json:"serviceName,omitempty"`
			ServiceCategory string `json:"serviceCategory,omitempty"`
		} `json:"serviceTypeDetail,omitempty"` //
		SmartPostInfoDetail struct {
			AncillaryEndorsement string `json:"ancillaryEndorsement,omitempty"`
			HubID                string `json:"hubId,omitempty"`
			Indicia              string `json:"indicia,omitempty"`
			SpecialServices      string `json:"specialServices,omitempty"`
		} `json:"smartPostInfoDetail,omitempty"` //
		ExpressFreightDetail struct {
			BookingConfirmationNumber string `json:"bookingConfirmationNumber,omitempty"`
			ShippersLoadAndCount      int    `json:"shippersLoadAndCount,omitempty"`
		} `json:"expressFreightDetail,omitempty"` //
		GroundShipment bool `json:"groundShipment,omitempty"` //
	} `json:"requestedShipment,omitempty"` //
	CarrierCodes []string `json:"carrierCodes,omitempty"` //
}

type RateResponse struct {
	TransactionID         string `json:"transactionId"`
	CustomerTransactionID string `json:"customerTransactionId"`
	Output                struct {
		RateReplyDetails []struct {
			ServiceType      string `json:"serviceType"`
			ServiceName      string `json:"serviceName"`
			PackagingType    string `json:"packagingType"`
			CustomerMessages []struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"customerMessages,omitempty"`
			RatedShipmentDetails []struct {
				RateType                         string  `json:"rateType"`
				RatedWeightMethod                string  `json:"ratedWeightMethod"`
				TotalDiscounts                   int     `json:"totalDiscounts"`
				TotalBaseCharge                  float64 `json:"totalBaseCharge"`
				TotalNetCharge                   float64 `json:"totalNetCharge"`
				TotalVatCharge                   int     `json:"totalVatCharge"`
				TotalNetFedExCharge              float64 `json:"totalNetFedExCharge"`
				TotalDutiesAndTaxes              int     `json:"totalDutiesAndTaxes"`
				TotalNetChargeWithDutiesAndTaxes float64 `json:"totalNetChargeWithDutiesAndTaxes"`
				TotalDutiesTaxesAndFees          int     `json:"totalDutiesTaxesAndFees"`
				TotalAncillaryFeesAndTaxes       int     `json:"totalAncillaryFeesAndTaxes"`
				ShipmentRateDetail               struct {
					RateZone             string  `json:"rateZone"`
					DimDivisor           int     `json:"dimDivisor"`
					FuelSurchargePercent float64 `json:"fuelSurchargePercent"`
					TotalSurcharges      float64 `json:"totalSurcharges"`
					TotalFreightDiscount int     `json:"totalFreightDiscount"`
					SurCharges           []struct {
						Type        string  `json:"type"`
						Description string  `json:"description"`
						Amount      float64 `json:"amount"`
					} `json:"surCharges"`
					PricingCode          string `json:"pricingCode"`
					CurrencyExchangeRate struct {
						FromCurrency string `json:"fromCurrency"`
						IntoCurrency string `json:"intoCurrency"`
						Rate         int    `json:"rate"`
					} `json:"currencyExchangeRate"`
					TotalBillingWeight struct {
						Units string `json:"units"`
						Valu  int    `json:"valu"`
					} `json:"totalBillingWeight"`
					Currency string `json:"currency"`
				} `json:"shipmentRateDetail,omitempty"`
				Currency string `json:"currency"`
			} `json:"ratedShipmentDetails,omitempty"`
			AnonymouslyAllowable bool `json:"anonymouslyAllowable,omitempty"`
			OperationalDetail    struct {
				OriginLocationIds                       string `json:"originLocationIds"`
				CommitDays                              string `json:"commitDays"`
				ServiceCode                             string `json:"serviceCode"`
				AirportID                               string `json:"airportId"`
				Scac                                    string `json:"scac"`
				OriginServiceAreas                      string `json:"originServiceAreas"`
				DeliveryDay                             string `json:"deliveryDay"`
				OriginLocationNumbers                   int    `json:"originLocationNumbers"`
				DestinationPostalCode                   string `json:"destinationPostalCode"`
				CommitDate                              string `json:"commitDate"`
				AstraDescription                        string `json:"astraDescription"`
				DeliveryDate                            string `json:"deliveryDate"`
				DeliveryEligibilities                   string `json:"deliveryEligibilities"`
				IneligibleForMoneyBackGuarantee         bool   `json:"ineligibleForMoneyBackGuarantee"`
				MaximumTransitTime                      string `json:"maximumTransitTime"`
				AstraPlannedServiceLevel                string `json:"astraPlannedServiceLevel"`
				DestinationLocationIds                  string `json:"destinationLocationIds"`
				DestinationLocationStateOrProvinceCodes string `json:"destinationLocationStateOrProvinceCodes"`
				TransitTime                             string `json:"transitTime"`
				PackagingCode                           string `json:"packagingCode"`
				DestinationLocationNumbers              int    `json:"destinationLocationNumbers"`
				PublishedDeliveryTime                   string `json:"publishedDeliveryTime"`
				CountryCodes                            string `json:"countryCodes"`
				StateOrProvinceCodes                    string `json:"stateOrProvinceCodes"`
				UrsaPrefixCode                          string `json:"ursaPrefixCode"`
				UrsaSuffixCode                          string `json:"ursaSuffixCode"`
				DestinationServiceAreas                 string `json:"destinationServiceAreas"`
				OriginPostalCodes                       string `json:"originPostalCodes"`
				CustomTransitTime                       string `json:"customTransitTime"`
			} `json:"operationalDetail,omitempty"`
			SignatureOptionType string `json:"signatureOptionType,omitempty"`
			ServiceDescription  struct {
				ServiceID   string `json:"serviceId"`
				ServiceType string `json:"serviceType"`
				Code        string `json:"code"`
				Names       []struct {
					Type     string `json:"type"`
					Encoding string `json:"encoding"`
					Value    string `json:"value"`
				} `json:"names"`
				OperatingOrgCodes []string `json:"operatingOrgCodes"`
				ServiceCategory   string   `json:"serviceCategory"`
				Description       string   `json:"description"`
				AstraDescription  string   `json:"astraDescription"`
			} `json:"serviceDescription,omitempty"`
			Commit struct {
				DateDetail struct {
					DayOfWeek    string `json:"dayOfWeek"`
					DayCxsFormat string `json:"dayCxsFormat"`
				} `json:"dateDetail"`
			} `json:"commit,omitempty"`
		} `json:"rateReplyDetails,omitempty"`
		QuoteDate string `json:"quoteDate"`
		Encoded   bool   `json:"encoded"`
		Alerts    []struct {
			Code      string `json:"code"`
			Message   string `json:"message"`
			AlertType string `json:"alertType"`
		} `json:"alerts"`
	} `json:"output,omitempty"`
	Errors []struct {
		Code    string `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
	} `json:"errors,omitempty"`
}

type RateXMLRequest struct {
	XMLName    xml.Name `xml:"SOAP-ENV:Envelope"`
	Xmlns_xsi  string   `xml:"xmlns xsi,attr,omitempty"`
	Xmlns_xsd  string   `xml:"xmlns xsd,attr,omitempty"`
	Xmlnsns    string   `xml:"xmlns xns,attr,omitempty"`
	Xmlns_soap string   `xml:"xmlns SOAP-ENV,attr,omitempty"`
	Xmlns_enc  string   `xml:"xmlns SOAP-ENC,attr,omitempty"`
	Text       string   `xml:",chardata"`
	Body       struct {
		Text        string `xml:",chardata"`
		RateRequest struct {
			Text                    string `xml:",chardata"`
			WebAuthenticationDetail struct {
				Text             string `xml:",chardata"`
				ParentCredential struct {
					Text     string `xml:",chardata"`
					Key      string `xml:"Key,omitempty"`
					Password string `xml:"Password,omitempty"`
				} `xml:"ParentCredential,omitempty"`
				UserCredential struct {
					Text     string `xml:",chardata"`
					Key      string `xml:"Key,omitempty"`
					Password string `xml:"Password,omitempty"`
				} `xml:"UserCredential"`
			} `xml:"WebAuthenticationDetail"`
			ClientDetail struct {
				Text          string `xml:",chardata"`
				AccountNumber string `xml:"AccountNumber,omitempty"`
				MeterNumber   string `xml:"MeterNumber,omitempty"`
				SoftwareId    string `xml:"SoftwareId,omitempty"`
			} `xml:"ClientDetail,omitempty"`
			TransactionDetail struct {
				Text                  string `xml:",chardata"`
				CustomerTransactionId string `xml:"CustomerTransactionId,omitempty"`
			} `xml:"TransactionDetail,omitempty"`
			Version struct {
				Text         string `xml:",chardata"`
				ServiceId    string `xml:"ServiceId,omitempty"`
				Major        string `xml:"Major,omitempty"`
				Intermediate string `xml:"Intermediate,omitempty"`
				Minor        string `xml:"Minor,omitempty"`
			} `xml:"Version,omitempty"`
			RequestedShipment struct {
				Text          string `xml:",chardata"`
				ShipTimestamp string `xml:"ShipTimestamp,omitempty"`
				DropoffType   string `xml:"DropoffType,omitempty"`
				ServiceType   string `xml:"ServiceType,omitempty"`
				PackagingType string `xml:"PackagingType,omitempty"`
				TotalWeight   struct {
					Text  string `xml:",chardata"`
					Units string `xml:"Units,omitempty"`
					Value string `xml:"Value,omitempty"`
				} `xml:"TotalWeight,omitempty"`
				Shipper struct {
					Text          string `xml:",chardata"`
					AccountNumber string `xml:"AccountNumber,omitempty"`
					Contact       struct {
						Text        string `xml:",chardata"`
						CompanyName string `xml:"CompanyName,omitempty"`
						PhoneNumber string `xml:"PhoneNumber,omitempty"`
					} `xml:"Contact,omitempty"`
					Address struct {
						Text                string   `xml:",chardata"`
						StreetLines         []string `xml:"StreetLines,omitempty"`
						City                string   `xml:"City,omitempty"`
						StateOrProvinceCode string   `xml:"StateOrProvinceCode,omitempty"`
						PostalCode          string   `xml:"PostalCode,omitempty"`
						CountryCode         string   `xml:"CountryCode,omitempty"`
					} `xml:"Address,omitempty"`
				} `xml:"Shipper,omitempty"`
				Recipient struct {
					Text          string `xml:",chardata"`
					AccountNumber string `xml:"AccountNumber,omitempty"`
					Contact       struct {
						Text        string `xml:",chardata"`
						PersonName  string `xml:"PersonName,omitempty"`
						PhoneNumber string `xml:"PhoneNumber,omitempty"`
					} `xml:"Contact,omitempty"`
					Address struct {
						Text                string   `xml:",chardata"`
						StreetLines         []string `xml:"StreetLines,omitempty"`
						City                string   `xml:"City,omitempty"`
						StateOrProvinceCode string   `xml:"StateOrProvinceCode,omitempty"`
						PostalCode          string   `xml:"PostalCode,omitempty"`
						CountryCode         string   `xml:"CountryCode,omitempty"`
						CountryName         string   `xml:"CountryName,omitempty"`
						Residential         string   `xml:"Residential,omitempty"`
					} `xml:"Address,omitempty"`
				} `xml:"Recipient,omitempty"`
				ShippingChargesPayment struct {
					Text        string `xml:",chardata"`
					PaymentType string `xml:"PaymentType,omitempty"`
					Payor       struct {
						Text             string `xml:",chardata"`
						ResponsibleParty struct {
							Text          string `xml:",chardata"`
							AccountNumber string `xml:"AccountNumber,omitempty"`
							Tins          struct {
								Text    string `xml:",chardata"`
								TinType string `xml:"TinType,omitempty"`
								Number  string `xml:"Number,omitempty"`
							} `xml:"Tins,omitempty"`
						} `xml:"ResponsibleParty,omitempty"`
					} `xml:"Payor,omitempty"`
				} `xml:"ShippingChargesPayment,omitempty"`
				RateRequestTypes          string `xml:"RateRequestTypes,omitempty"`
				PackageCount              string `xml:"PackageCount,omitempty"`
				RequestedPackageLineItems struct {
					Text              string `xml:",chardata"`
					SequenceNumber    string `xml:"SequenceNumber,omitempty"`
					GroupNumber       string `xml:"GroupNumber,omitempty"`
					GroupPackageCount string `xml:"GroupPackageCount,omitempty"`
					Weight            struct {
						Text  string `xml:",chardata"`
						Units string `xml:"Units,omitempty"`
						Value string `xml:"Value,omitempty"`
					} `xml:"Weight,omitempty"`
					Dimensions struct {
						Text   string `xml:",chardata"`
						Length string `xml:"Length,omitempty"`
						Width  string `xml:"Width,omitempty"`
						Height string `xml:"Height,omitempty"`
						Units  string `xml:"Units,omitempty"`
					} `xml:"Dimensions,omitempty"`
					ContentRecords struct {
						Text             string `xml:",chardata"`
						PartNumber       string `xml:"PartNumber,omitempty"`
						ItemNumber       string `xml:"ItemNumber,omitempty"`
						ReceivedQuantity string `xml:"ReceivedQuantity,omitempty"`
						Description      string `xml:"Description,omitempty"`
					} `xml:"ContentRecords,omitempty"`
				} `xml:"RequestedPackageLineItems,omitempty"`
			} `xml:"RequestedShipment,omitempty"`
		} `xml:"RateRequest,omitempty"`
	} `xml:"SOAP-ENV:Body"`
}

type RateXMLResponse struct {
	Body struct {
		Text  string `xml:",chardata"`
		Fault struct {
			Text        string `xml:",chardata"`
			Faultcode   string `xml:"faultcode,omitempty"`
			Faultstring struct {
				Text string `xml:",chardata"`
				Lang string `xml:"lang,attr,omitempty"`
			} `xml:"faultstring,omitempty"`
			Detail struct {
				Text  string `xml:",chardata"`
				Cause string `xml:"cause,omitempty"`
				Code  string `xml:"code,omitempty"`
				Desc  string `xml:"desc,omitempty"`
			} `xml:"detail,omitempty"`
		} `xml:"Fault,omitempty"`
		RateReply struct {
			Text            string `xml:",chardata"`
			Xmlns           string `xml:"xmlns,attr"`
			HighestSeverity string `xml:"HighestSeverity"`
			Notifications   struct {
				Text             string `xml:",chardata"`
				Severity         string `xml:"Severity,omitempty"`
				Source           string `xml:"Source,omitempty"`
				Code             string `xml:"Code,omitempty"`
				Message          string `xml:"Message,omitempty"`
				LocalizedMessage string `xml:"LocalizedMessage,omitempty"`
			} `xml:"Notifications,omitempty"`
			TransactionDetail struct {
				Text                  string `xml:",chardata"`
				CustomerTransactionId string `xml:"CustomerTransactionId,omitempty"`
			} `xml:"TransactionDetail,omitempty"`
			Version struct {
				Text         string `xml:",chardata"`
				ServiceId    string `xml:"ServiceId,omitempty"`
				Major        string `xml:"Major,omitempty"`
				Intermediate string `xml:"Intermediate,omitempty"`
				Minor        string `xml:"Minor,omitempty"`
			} `xml:"Version,omitempty"`
			RateReplyDetails struct {
				Text               string `xml:",chardata"`
				ServiceType        string `xml:"ServiceType,omitempty"`
				ServiceDescription struct {
					Text        string `xml:",chardata"`
					ServiceType string `xml:"ServiceType,omitempty"`
					Code        string `xml:"Code,omitempty"`
					Names       []struct {
						Text     string `xml:",chardata"`
						Type     string `xml:"Type,omitempty"`
						Encoding string `xml:"Encoding,omitempty"`
						Value    string `xml:"Value,omitempty"`
					} `xml:"Names,omitempty"`
					Description      string `xml:"Description,omitempty"`
					AstraDescription string `xml:"AstraDescription,omitempty"`
				} `xml:"ServiceDescription,omitempty"`
				PackagingType                   string `xml:"PackagingType,omitempty"`
				DestinationAirportId            string `xml:"DestinationAirportId,omitempty"`
				IneligibleForMoneyBackGuarantee string `xml:"IneligibleForMoneyBackGuarantee,omitempty"`
				SignatureOption                 string `xml:"SignatureOption,omitempty"`
				ActualRateType                  string `xml:"ActualRateType,omitempty"`
				RatedShipmentDetails            []struct {
					Text                 string `xml:",chardata"`
					EffectiveNetDiscount struct {
						Text     string `xml:",chardata"`
						Currency string `xml:"Currency,omitempty"`
						Amount   string `xml:"Amount,omitempty"`
					} `xml:"EffectiveNetDiscount,omitempty"`
					ShipmentRateDetail struct {
						Text                 string `xml:",chardata"`
						RateType             string `xml:"RateType,omitempty"`
						RateZone             string `xml:"RateZone,omitempty"`
						RatedWeightMethod    string `xml:"RatedWeightMethod,omitempty"`
						DimDivisor           string `xml:"DimDivisor,omitempty"`
						FuelSurchargePercent string `xml:"FuelSurchargePercent,omitempty"`
						TotalBillingWeight   struct {
							Text  string `xml:",chardata"`
							Units string `xml:"Units,omitempty"`
							Value string `xml:"Value,omitempty"`
						} `xml:"TotalBillingWeight,omitempty"`
						TotalBaseCharge struct {
							Text     string `xml:",chardata"`
							Currency string `xml:"Currency,omitempty"`
							Amount   string `xml:"Amount,omitempty"`
						} `xml:"TotalBaseCharge,omitempty"`
						TotalFreightDiscounts struct {
							Text     string `xml:",chardata"`
							Currency string `xml:"Currency,omitempty"`
							Amount   string `xml:"Amount,omitempty"`
						} `xml:"TotalFreightDiscounts,omitempty"`
						TotalNetFreight struct {
							Text     string `xml:",chardata"`
							Currency string `xml:"Currency,omitempty"`
							Amount   string `xml:"Amount,omitempty"`
						} `xml:"TotalNetFreight,omitempty"`
						TotalSurcharges struct {
							Text     string `xml:",chardata"`
							Currency string `xml:"Currency,omitempty"`
							Amount   string `xml:"Amount,omitempty"`
						} `xml:"TotalSurcharges,omitempty"`
						TotalNetFedExCharge struct {
							Text     string `xml:",chardata"`
							Currency string `xml:"Currency,omitempty"`
							Amount   string `xml:"Amount,omitempty"`
						} `xml:"TotalNetFedExCharge,omitempty"`
						TotalTaxes struct {
							Text     string `xml:",chardata"`
							Currency string `xml:"Currency,omitempty"`
							Amount   string `xml:"Amount,omitempty"`
						} `xml:"TotalTaxes,omitempty"`
						TotalNetCharge struct {
							Text     string `xml:",chardata"`
							Currency string `xml:"Currency,omitempty"`
							Amount   string `xml:"Amount,omitempty"`
						} `xml:"TotalNetCharge,omitempty"`
						TotalRebates struct {
							Text     string `xml:",chardata"`
							Currency string `xml:"Currency,omitempty"`
							Amount   string `xml:"Amount,omitempty"`
						} `xml:"TotalRebates,omitempty"`
						TotalDutiesAndTaxes struct {
							Text     string `xml:",chardata"`
							Currency string `xml:"Currency,omitempty"`
							Amount   string `xml:"Amount,omitempty"`
						} `xml:"TotalDutiesAndTaxes,omitempty"`
						TotalAncillaryFeesAndTaxes struct {
							Text     string `xml:",chardata"`
							Currency string `xml:"Currency,omitempty"`
							Amount   string `xml:"Amount,omitempty"`
						} `xml:"TotalAncillaryFeesAndTaxes,omitempty"`
						TotalDutiesTaxesAndFees struct {
							Text     string `xml:",chardata"`
							Currency string `xml:"Currency,omitempty"`
							Amount   string `xml:"Amount,omitempty"`
						} `xml:"TotalDutiesTaxesAndFees,omitempty"`
						TotalNetChargeWithDutiesAndTaxes struct {
							Text     string `xml:",chardata"`
							Currency string `xml:"Currency,omitempty"`
							Amount   string `xml:"Amount,omitempty"`
						} `xml:"TotalNetChargeWithDutiesAndTaxes,omitempty"`
						Surcharges []struct {
							Text          string `xml:",chardata"`
							SurchargeType string `xml:"SurchargeType,omitempty"`
							Level         string `xml:"Level,omitempty"`
							Description   string `xml:"Description,omitempty"`
							Amount        struct {
								Text     string `xml:",chardata"`
								Currency string `xml:"Currency,omitempty"`
								Amount   string `xml:"Amount,omitempty"`
							} `xml:"Amount,omitempty"`
						} `xml:"Surcharges,omitempty"`
					} `xml:"ShipmentRateDetail,omitempty"`
					RatedPackages struct {
						Text                 string `xml:",chardata"`
						GroupNumber          string `xml:"GroupNumber,omitempty"`
						EffectiveNetDiscount struct {
							Text     string `xml:",chardata"`
							Currency string `xml:"Currency,omitempty"`
							Amount   string `xml:"Amount,omitempty"`
						} `xml:"EffectiveNetDiscount,omitempty"`
						PackageRateDetail struct {
							Text              string `xml:",chardata"`
							RateType          string `xml:"RateType,omitempty"`
							RatedWeightMethod string `xml:"RatedWeightMethod,omitempty"`
							BillingWeight     struct {
								Text  string `xml:",chardata"`
								Units string `xml:"Units,omitempty"`
								Value string `xml:"Value,omitempty"`
							} `xml:"BillingWeight,omitempty"`
							BaseCharge struct {
								Text     string `xml:",chardata"`
								Currency string `xml:"Currency,omitempty"`
								Amount   string `xml:"Amount,omitempty"`
							} `xml:"BaseCharge,omitempty"`
							TotalFreightDiscounts struct {
								Text     string `xml:",chardata"`
								Currency string `xml:"Currency,omitempty"`
								Amount   string `xml:"Amount,omitempty"`
							} `xml:"TotalFreightDiscounts,omitempty"`
							NetFreight struct {
								Text     string `xml:",chardata"`
								Currency string `xml:"Currency,omitempty"`
								Amount   string `xml:"Amount,omitempty"`
							} `xml:"NetFreight,omitempty"`
							TotalSurcharges struct {
								Text     string `xml:",chardata"`
								Currency string `xml:"Currency,omitempty"`
								Amount   string `xml:"Amount,omitempty"`
							} `xml:"TotalSurcharges,omitempty"`
							NetFedExCharge struct {
								Text     string `xml:",chardata"`
								Currency string `xml:"Currency,omitempty"`
								Amount   string `xml:"Amount,omitempty"`
							} `xml:"NetFedExCharge,omitempty"`
							TotalTaxes struct {
								Text     string `xml:",chardata"`
								Currency string `xml:"Currency,omitempty"`
								Amount   string `xml:"Amount,omitempty"`
							} `xml:"TotalTaxes,omitempty"`
							NetCharge struct {
								Text     string `xml:",chardata"`
								Currency string `xml:"Currency,omitempty"`
								Amount   string `xml:"Amount,omitempty"`
							} `xml:"NetCharge,omitempty"`
							TotalRebates struct {
								Text     string `xml:",chardata"`
								Currency string `xml:"Currency,omitempty"`
								Amount   string `xml:"Amount,omitempty"`
							} `xml:"TotalRebates,omitempty"`
							Surcharges []struct {
								Text          string `xml:",chardata"`
								SurchargeType string `xml:"SurchargeType,omitempty"`
								Level         string `xml:"Level,omitempty"`
								Description   string `xml:"Description,omitempty"`
								Amount        struct {
									Text     string `xml:",chardata"`
									Currency string `xml:"Currency,omitempty"`
									Amount   string `xml:"Amount,omitempty"`
								} `xml:"Amount,omitempty"`
							} `xml:"Surcharges,omitempty"`
						} `xml:"PackageRateDetail,omitempty"`
					} `xml:"RatedPackages,omitempty"`
				} `xml:"RatedShipmentDetails,omitempty"`
			} `xml:"RateReplyDetails,omitempty"`
		} `xml:"RateReply,omitempty"`
	} `xml:"Body"`
}

func (c RateRequest) Rate(token string, apiUrl string) (RateResponse, error) {
	var _response RateResponse
	client := &http.Client{}

	reqUrl := apiUrl + "/rate/v1/rates/quotes"

	request, errJSON := json.Marshal(c)

	if errJSON != nil {
		return _response, errJSON
	}

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(request))
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("X-locale", "en_US")

	resp, err := client.Do(req)

	if err != nil {
		return _response, err
	}
	defer resp.Body.Close()

	log.Println(resp.StatusCode)

	if resp.StatusCode == 200 {
		data := json.NewDecoder(resp.Body)
		errjson := data.Decode(&_response)
		if errjson != nil {
			log.Println(errjson)
		}
		return _response, nil
	} else {
		dataError := json.NewDecoder(resp.Body)
		errjson := dataError.Decode(&_response)
		if errjson != nil {
			log.Println(errjson)
		}
		return _response, errors.New(resp.Status)
	}

	return _response, nil
}

func (c RateXMLRequest) Rate(url string, testMode bool) (RateXMLResponse, error) {
	var _response RateXMLResponse
	request, _ := xml.Marshal(c)
	newStr := `SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns="http://fedex.com/ws/rate/v28"`
	s := strings.Replace(string(request), "SOAP-ENV:Envelope", newStr, 1)

	content, err, statusCode := common.Fedex{TestMode: testMode}.PostRequest(s, url)

	if err != nil {
		log.Printf("%s", err)
		return RateXMLResponse{}, err
	}

	if statusCode == 503 {
		_statusCode, _ := conv.String(statusCode)
		return RateXMLResponse{}, errors.New("Backend Error with code " + _statusCode)
	}

	log.Printf("%s", content)

	err = xml.Unmarshal(content, &_response)
	if err != nil {
		log.Println("%s", err)
		return RateXMLResponse{}, err
	}

	//log.Printf("%s", _response.Body)
	return _response, nil
}
