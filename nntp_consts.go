package nntp

import "strconv"

type ResponseCode int

func (code ResponseCode) Error() string {
	return strconv.Itoa(int(code))
}

const (
	// Connection

	/**
	 * 'Server ready - posting allowed' (RFC977)
	 */
	ResponseCodeReadyPostingAllowed ResponseCode = 200

	/**
	 * 'Server ready - no posting allowed' (RFC977)
	 */
	ResponseCodeReadyPostingProhibited ResponseCode = 201

	/**
	 * 'Closing connection - goodbye!' (RFC977)
	 */
	ResponseCodeDisconnectingRequested ResponseCode = 205 ///// goodbye

	/**
	 * 'Service discontinued' (RFC977)
	 */
	ResponseCodeDisconnectingForced ResponseCode = 400 ///// unavailable / discontinued

	/**
	 * 'Slave status noted' (RFC977)
	 */
	ResponseCodeSlaveRecognized ResponseCode = 202

	// Common errors

	/**
	 * 'Command not recognized' (RFC977)
	 */
	ResponseCodeUnknownCommand ResponseCode = 500

	/**
	 * 'Command syntax error' (RFC977)
	 */
	ResponseCodeSyntaxError ResponseCode = 501

	/**
	 * 'Access restriction or permission denied' (RFC977)
	 */
	ResponseCodeNotPermitted ResponseCode = 502

	/**
	 * 'Program fault - command not performed' (RFC977)
	 */
	ResponseCodeNotSupported ResponseCode = 503

	// Group selection

	/**
	 * 'Group selected' (RFC977)
	 */
	ResponseCodeGroupSelected ResponseCode = 211

	/**
	 * 'No such news group' (RFC977)
	 */
	ResponseCodeNoSuchGroup ResponseCode = 411

	// Article retrieval

	/**
	 * 'Article retrieved - head and body follow' (RFC977)
	 */
	ResponseCodeArticleFollows ResponseCode = 220

	/**
	 * 'Article retrieved - head follows' (RFC977)
	 */
	ResponseCodeHeadFollows ResponseCode = 221

	/**
	 * 'Article retrieved - body follows' (RFC977)
	 */
	ResponseCodeBodyFollows ResponseCode = 222

	/**
	 * 'Article retrieved - request text separately' (RFC977)
	 */
	ResponseCodeArticleSelected ResponseCode = 223

	/**
	 * 'No newsgroup has been selected' (RFC977)
	 */
	ResponseCodeNoGroupSelected ResponseCode = 412

	/**
	 * 'No current article has been selected' (RFC977)
	 */
	ResponseCodeNoArticleSelected ResponseCode = 420

	/**
	 * 'No next article in this group' (RFC977)
	 */
	ResponseCodeNoNextArticle ResponseCode = 421

	/**
	 * 'No previous article in this group' (RFC977)
	 */
	ResponseCodeNoPreviousArticle ResponseCode = 422

	/**
	 * 'No such article number in this group' (RFC977)
	 */
	ResponseCodeNoSuchArticleNumber ResponseCode = 423

	/**
	 * 'No such article found' (RFC977)
	 */
	ResponseCodeNoSuchArticleId ResponseCode = 430

	// Transferring

	/**
	 * 'Send article to be transferred' (RFC977)
	 */
	ResponseCodeTransferSend ResponseCode = 335

	/**
	 * 'Article transferred ok' (RFC977)
	 */
	ResponseCodeTransferSuccess ResponseCode = 235

	/**
	 * 'Article not wanted - do not send it' (RFC977)
	 */
	ResponseCodeTransferUnwanted ResponseCode = 435

	/**
	 * 'Transfer failed - try again later' (RFC977)
	 */
	ResponseCodeTransferFailure ResponseCode = 436

	/**
	 * 'Article rejected - do not try again' (RFC977)
	 */
	ResponseCodeTransferRejected ResponseCode = 437

	// Posting

	/**
	 * 'Send article to be posted' (RFC977)
	 */
	ResponseCodePostingSend ResponseCode = 340

	/**
	 * 'Article posted ok' (RFC977)
	 */
	ResponseCodePostingSuccess ResponseCode = 240

	/**
	 * 'Posting not allowed' (RFC977)
	 */
	ResponseCodePostingProhibited ResponseCode = 440

	/**
	 * 'Posting failed' (RFC977)
	 */
	ResponseCodePostingFailure ResponseCode = 441

	// Authorization

	/**
	 * 'Authorization required for this command' (RFC2980)
	 */
	ResponseCodeAuthorizationRequired ResponseCode = 450

	/**
	 * 'Continue with authorization sequence' (RFC2980)
	 */
	ResponseCodeAuthorizationContinue ResponseCode = 350

	/**
	 * 'Authorization accepted' (RFC2980)
	 */
	ResponseCodeAuthorizationAccepted ResponseCode = 250

	/**
	 * 'Authorization rejected' (RFC2980)
	 */
	ResponseCodeAuthorizationRejected ResponseCode = 452

	// Authentication

	/**
	 * 'Authentication required' (RFC2980)
	 */
	ResponseCodeAuthenticationRequired ResponseCode = 480

	/**
	 * 'More authentication information required' (RFC2980)
	 */
	ResponseCodeAuthenticationContinue ResponseCode = 381

	/**
	 * 'Authentication accepted' (RFC2980)
	 */
	ResponseCodeAuthenticationAccepted ResponseCode = 281

	/**
	 * 'Authentication rejected' (RFC2980)
	 */
	ResponseCodeAuthenticationRejected ResponseCode = 482

	// Misc

	/**
	 * 'Help text follows' (Draft)
	 */
	ResponseCodeHelpFollows ResponseCode = 100

	/**
	 * 'Capabilities list follows' (Draft)
	 */
	ResponseCodeCapabilitiesFollow ResponseCode = 101

	/**
	 * 'Server date and time' (Draft)
	 */
	ResponseCodeServerDate ResponseCode = 111

	/**
	 * 'Information follows' (Draft)
	 */
	ResponseCodeInformationFollows ResponseCode = 215

	/**
	 * 'Groups follows' (Draft)
	 */
	ResponseCodeGroupsFollow ResponseCode = 215

	/**
	 * 'Overview information follows' (Draft)
	 */
	ResponseCodeOverviewFollows ResponseCode = 224

	/**
	 * 'Headers follow' (Draft)
	 */
	ResponseCodeHeadersFollow ResponseCode = 225

	/**
	 * 'List of new articles follows' (Draft)
	 */
	ResponseCodeNewArticlesFollow ResponseCode = 230

	/**
	 * 'List of new newsgroups follows' (Draft)
	 */
	ResponseCodeNewGroupsFollow ResponseCode = 231

	/**
	 * 'The server is in the wrong mode; the indicated capability should be used to change the mode' (Draft)
	 */
	ResponseCodeWrongMode ResponseCode = 401

	/**
	 * 'Internal fault or problem preventing action being taken' (Draft)
	 */
	ResponseCodeInternalFault ResponseCode = 403

	/**
	 * 'Command unavailable until suitable privacy has been arranged' (Draft)
	 *
	 * (the client must negotiate appropriate privacy protection on the connection.
	 * This will involve the use of a privacy extension such as [NNTP-TLS].)
	 */
	ResponseCodeEncryptionRequired ResponseCode = 483

	/**
	 * 'Error in base64-encoding [RFC3548] of an argument' (Draft)
	 */
	ResponseCodeBase64EncodingError ResponseCode = 504
)

type GroupPermission byte

const (
	GroupPostingPermitted GroupPermission = 'y'
	GroupPostingForbidden GroupPermission = 'n'
	GroupPostingModerated GroupPermission = 'm'
)

const DefaultArticleDateLayout = "02 Jan 2006 15:04:05 -0700"
const AlternativeArticleDateLayout = "02 Jan 2006 15:04:05 MST"

type OverviewFieldType int

const (
	ShortHeaderOverviewField OverviewFieldType = iota
	FullHeaderOverviewField
	MetadataOverviewField
)
