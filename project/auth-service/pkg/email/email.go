package email

import (
	"auth-service/pkg/custom_errors"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"html/template"
)

const message = `{{define "authCode"}}
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>TechSaler</title>
	<style type="text/css">
		html {
			-webkit-text-size-adjust: none;
			-ms-text-size-adjust: none;
		}
		@media only screen and (max-device-width:660px),
		only screen and (max-width:660px) {
			.em-mob-text_align-left {
				text-align: left !important;
			}
			.em-mob-width-100perc {
				width: 100% !important;
				max-width: 100% !important;
			}
			.em-mob-text_align-center {
				text-align: center !important;
			}
			.em-mob-font_size-14px {
				font-size: 14px !important;
				line-height: 20px !important;
			}
			.em-mob-width-91perc {
				width: 91% !important;
				max-width: 91% !important;
			}
		}
	</style>
</head>
<body style="margin: 0; padding: 0;">
	<span class="preheader" style="display: none !important; visibility: hidden; opacity: 0; color: #F8F8F8; height: 0; width: 0; font-size: 1px;">
		Прехедер
		&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;
	</span>
	<!--[if !mso]><!-->
	<div style="font-size:0px;color:transparent;opacity:0;">
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	</div>
	<!--<![endif]-->
	<table cellpadding="0" cellspacing="0" border="0" width="100%" style="font-size: 1px; line-height: normal;">
		<tr em="group">
			<td align="center" bgcolor="#F8F8F8">
				<!--[if (gte mso 9)|(IE)]>
				<table cellpadding="0" cellspacing="0" border="0" width="800"><tr><td>
				<![endif]-->
				<table cellpadding="0" cellspacing="0" width="100%" border="0" style="max-width: 800px; min-width: 320px; width: 100%;">
					<tr em="block">
						<td align="center">
							<!--[if (gte mso 9)|(IE)]>
							<table cellpadding="0" cellspacing="0" border="0" width="660"><tr><td>
							<![endif]-->
							<table align="center" cellpadding="0" cellspacing="0" border="0" width="100%" style="max-width: 660px;">
								<tr>
									<td align="center" bgcolor="#FFFFFF">
										<table cellpadding="0" cellspacing="0" border="0" width="100%">
											<tr>
												<td height="15"></td>
											</tr>
										</table>
										<table align="center" cellpadding="0" cellspacing="0" border="0" width="91%">
											<tr><td align="left" valign="top">
													<div style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 12px; line-height: 14px; color: #9299a2;"> Осталось подтвердить e-mail</div>
												</td><td width="20">&nbsp;</td></tr>
										</table>
										<table cellpadding="0" cellspacing="0" border="0" width="100%">
											<tr>
												<td height="10"></td>
											</tr>
										</table>
									</td>
								</tr>
							</table>
							<!--[if (gte mso 9)|(IE)]>
							</td></tr></table>
							<![endif]-->
						</td>
					</tr>
					<tr em="block">
						<td align="center">
							<!--[if (gte mso 9)|(IE)]>
							<table cellpadding="0" cellspacing="0" border="0" width="660"><tr><td>
							<![endif]-->
							<table align="center" width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 660px;">
								<tr>
									<td align="center" valign="middle" bgcolor="#FFFFFF">
										<table cellpadding="0" cellspacing="0" border="0" width="100%">
											<tr>
												<td height="20"></td>
											</tr>
										</table>
										<!--[if (gte mso 9)|(IE)]>
										<table border="0" cellspacing="0" cellpadding="0">
										<tr><td width="260"><![endif]-->
										<div style="display: inline-block; width: 260px; vertical-align: middle;">
											<table width="100%" border="0" cellspacing="0" cellpadding="0">
												<tr>
													<td align="left" class="em-mob-text_align-center">
														<a href="https://bytrip.ru/" target="_blank"><img src="https://emcdn.ru/169769/230721_2791_GpoFxhb.png" border="0" alt="" style="display: inline-block; width: 100%; max-width: 190px;" width="190"></a>
														<table cellpadding="0" cellspacing="0" border="0" width="100%">
															<tr>
																<td height="26"></td>
															</tr>
														</table>
													</td>
												</tr>
											</table>
										</div>
										<!--[if (gte mso 9)|(IE)]></td><td width="340"><![endif]-->
										<div style="display: inline-block; width: 340px; vertical-align: middle;" class="em-mob-width-91perc">
											<table width="100%" border="0" cellspacing="0" cellpadding="0">
												<tr>
													<td align="right">
														<table align="center" width="100%" border="0" cellspacing="0" cellpadding="0">
															<tr>
																<td align="center" em="atom">
																	<a style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 16px; line-height: 24px; color: #333333; text-decoration: none; white-space: nowrap;" class="em-mob-font_size-14px" target="_blank">&nbsp;</a>
																</td>
																<td width="15" em="atom">&nbsp;</td>
																<td align="center" em="atom">
																	<a style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 16px; line-height: 24px; color: #333333; text-decoration: none; white-space: nowrap;" class="em-mob-font_size-14px" target="_blank">&nbsp;</a>
																</td>
																<td width="15" em="atom">&nbsp;</td>
																<td align="center" em="atom">
																	<a style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 16px; line-height: 24px; color: #333333; text-decoration: none; white-space: nowrap;" class="em-mob-font_size-14px" target="_blank" href="hotel.bytrip.ru">Для партнёров&nbsp;</a>
																</td><td align="center" em="atom">
																	<a style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 16px; line-height: 24px; color: #333333; text-decoration: none; white-space: nowrap;" class="em-mob-font_size-14px" target="_blank" href="https://bytrip.ru/">Выбрать отель&nbsp;</a>
																</td>
																<td width="15" em="atom">&nbsp;</td>
																<td align="center" em="atom">
																	<a href="https://bytrip.ru/" target="_blank" style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 16px; line-height: 24px; color: #276eaf; text-decoration: none; white-space: nowrap;" class="em-mob-font_size-14px">Личный кабинет</a>
																</td>
															</tr>
														</table>
														<table cellpadding="0" cellspacing="0" border="0" width="100%">
															<tr>
																<td height="20"></td>
															</tr>
														</table>
													</td>
												</tr>
											</table>
										</div>
										<!--[if (gte mso 9)|(IE)]>
										</td></tr>
										</table><![endif]-->
										<table cellpadding="0" cellspacing="0" border="0" width="100%">
											<tr>
												<td height="10"></td>
											</tr>
										</table>
									</td>
								</tr>
							</table>
							<!--[if (gte mso 9)|(IE)]>
							</td></tr></table>
							<![endif]-->
						</td>
					</tr>
					<tr em="block">
						<td align="center">
							<!--[if (gte mso 9)|(IE)]>
							<table cellpadding="0" cellspacing="0" border="0" width="660"><tr><td>
							<![endif]-->
							<table align="center" cellpadding="0" cellspacing="0" border="0" width="100%" style="max-width: 660px;">
								<tr>
									<td align="center" bgcolor="#FFFFFF">
										<table align="center" cellpadding="0" cellspacing="0" border="0" width="91%">
											<tr>
												<td align="left">
													<table cellpadding="0" cellspacing="0" border="0" width="100%">
														<tr>
															<td height="20"></td>
														</tr>
													</table>
													<div style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 36px; line-height: 44px; color: #333333;"> <strong>Подтвердите адрес электронной почты</strong></div>
													<table cellpadding="0" cellspacing="0" border="0" width="100%">
														<tr>
															<td height="10"></td>
														</tr>
													</table>
												</td>
											</tr>
										</table>
									</td>
								</tr>
							</table>
							<!--[if (gte mso 9)|(IE)]>
							</td></tr></table>
							<![endif]-->
						</td>
					</tr>
					<tr em="block">
						<td align="center">
							<!--[if (gte mso 9)|(IE)]>
							<table cellpadding="0" cellspacing="0" border="0" width="660"><tr><td>
							<![endif]-->
							<table align="center" cellpadding="0" cellspacing="0" border="0" width="100%" style="max-width: 660px;">
								<tr>
									<td align="center" bgcolor="#FFFFFF">
										<table align="center" cellpadding="0" cellspacing="0" border="0" width="91%">
											<tr>
												<td align="left">
													<table cellpadding="0" cellspacing="0" border="0" width="100%">
														<tr>
															<td height="20"></td>
														</tr>
													</table>
													<table cellpadding="0" cellspacing="0" border="0" width="100%">
														<tr></tr>
													</table>
													<table cellpadding="0" cellspacing="0" border="0" width="100%">
														<tr></tr>
													</table>
													<div style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 16px; line-height: 24px; color: #333333;">&nbsp;Для подтверждения почты и продолжения регистрации, введите код ниже в соответствующее окно</div>
													<table cellpadding="0" cellspacing="0" border="0" width="100%">
														<tr>
															<td height="20"></td>
														</tr>
													</table>
												</td>
											</tr>
										</table>
									</td>
								</tr>
							</table>
							<!--[if (gte mso 9)|(IE)]>
							</td></tr></table>
							<![endif]-->
						</td>
					</tr>
					<tr em="block">
						<td align="center">
							<!--[if (gte mso 9)|(IE)]>
							<table cellpadding="0" cellspacing="0" border="0" width="660"><tr><td>
							<![endif]-->
							<table align="center" cellpadding="0" cellspacing="0" border="0" width="100%" style="max-width: 660px;">
								<tr>
									<td align="center" bgcolor="#FFFFFF">
										<table align="center" cellpadding="0" cellspacing="0" border="0" width="91%">
											<tr>
												<td align="center">
													<table cellpadding="0" cellspacing="0" border="0" width="100%">
														<tr>
															<td height="20"></td>
														</tr>
													</table>
													<table cellpadding="0" cellspacing="0" border="0" width="290">
														<tr>
															<td align="center" valign="middle" height="45" style="background-color: #276EAF; border-radius: 7px; height: 45px;">
																<a href="" target="_blank" style="display: block; width: 100%; height: 45px; font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; color: #ffffff; font-size: 16px; line-height: 45px; text-decoration: none; white-space: nowrap;">{{.}}<br></a>
															</td>
														</tr>
													</table>
													<table cellpadding="0" cellspacing="0" border="0" width="100%">
														<tr>
															<td height="50"></td>
														</tr>
													</table>
												</td>
											</tr>
										</table>
									</td>
								</tr>
							</table>
							<!--[if (gte mso 9)|(IE)]>
							</td></tr></table>
							<![endif]-->
						</td>
					</tr>
					<tr em="block">
						<td align="center">
							<!--[if (gte mso 9)|(IE)]>
							<table cellpadding="0" cellspacing="0" border="0" width="660"><tr><td>
							<![endif]-->
							<table align="center" width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 660px;">
								<tr>
									<td align="center" valign="top">
									</td>
								</tr>
							</table>
							<!--[if (gte mso 9)|(IE)]>
							</td></tr></table>
							<![endif]-->
						</td>
					</tr>
					
					<tr em="block">
						<td align="center">
							<table cellpadding="0" cellspacing="0" border="0" width="100%">
								<tr>
									<td height="30"></td>
								</tr>
							</table>
							<table align="center" width="100%" border="0" cellspacing="0" cellpadding="0">
								<tr>
									<td align="center">
										<!--[if (gte mso 9)|(IE)]>
										<table border="0" cellspacing="0" cellpadding="0">
										<tr><td width="310"><![endif]-->
										<div style="display: inline-block; width: 310px; vertical-align: top;" class="em-mob-width-100perc">
											<table width="290" border="0" cellspacing="0" cellpadding="0" class="em-mob-width-91perc">
												<tr>
													<td align="left">
														<a href="https://bytrip.ru/" target="_blank"><img src="https://emcdn.ru/169769/230721_2791_RLDpdEN.png" width="190" border="0" alt="" style="display: inline-block; width: 100%; max-width: 190px;"></a>
														<table cellpadding="0" cellspacing="0" border="0" width="100%">
															<tr>
																<td height="25"></td>
															</tr>
														</table>
													</td>
												</tr>
											</table>
										</div>
										<!--[if (gte mso 9)|(IE)]></td><td width="310"><![endif]-->
										<div style="display: inline-block; width: 310px; vertical-align: top;" class="em-mob-width-100perc">
											<table width="290" border="0" cellspacing="0" cellpadding="0" class="em-mob-width-91perc">
												<tr>
													<td align="right" class="em-mob-text_align-left">
														<table border="0" cellspacing="0" cellpadding="0">
															<tr>
																<td>
																	<a href="" target="_blank"></a>
																</td>
																<td width="10">&nbsp;</td>
																<td>
																	<a href="" target="_blank"></a>
																</td>
																<td width="10">&nbsp;</td>
																<td>
																	<a href="" target="_blank"></a>
																</td>
																<td width="10">&nbsp;</td>
																<td>
																	<a href="" target="_blank"></a>
																</td>
																<td width="10">&nbsp;</td>
																<td>
																	<a href="" target="_blank"></a>
																</td>
																<td width="10">&nbsp;</td>
																<td>
																	<a href="" target="_blank"></a>
																</td>
																<td width="10">&nbsp;</td>
																<td>
																	<a href="" target="_blank"></a>
																</td>
															</tr>
														</table>
														<table cellpadding="0" cellspacing="0" border="0" width="100%">
															<tr>
																<td height="25"></td>
															</tr>
														</table>
													</td>
												</tr>
											</table>
										</div>
										<!--[if (gte mso 9)|(IE)]>
										</td></tr>
										</table><![endif]-->
									</td>
								</tr>
							</table>
							<!--[if (gte mso 9)|(IE)]>
							<table cellpadding="0" cellspacing="0" border="0" width="660"><tr><td>
							<![endif]-->
							<table align="center" width="91%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px;">
								<tr>
									<td align="left" valign="top">
										<table border="0" cellspacing="0" cellpadding="0">
											<tr><td>
													<a href="https://t.me/bytripru" target="_blank"><img src="https://imgems.ru/emailmaker/techsalerlight/tg-2.png" width="30" border="0" alt="" style="display: block;"></a>
												</td><td width="10">&nbsp;</td><td>
													
												</td><td>
													<a href="https://vk.com/bytripru" target="_blank"><img src="https://imgems.ru/emailmaker/techsalerlight/vk-2.png" width="30" border="0" alt="" style="display: block; max-width: 30px;"></a>
												</td><td width="10">&nbsp;</td><td>
													
												</td><td width="10">&nbsp;</td><td>
													
												</td><td width="10">&nbsp;</td><td>
													
												</td><td width="10">&nbsp;</td><td>
													
												</td></tr>
										</table>
										<table align="center" width="100%" border="0" cellspacing="0" cellpadding="0">
											<tr>
											</tr>
										</table>
										<table border="0" cellspacing="0" cellpadding="0">
											<tr><td align="center">
													<a style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 14px; line-height: 21px; color: #333333; text-decoration: none;" target="_blank">&nbsp;</a>
												</td><td align="center">
													<a style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 14px; line-height: 21px; color: #333333; text-decoration: none;" target="_blank">&nbsp;</a>
												</td><td align="center">
													<a href="" target="_blank" style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 14px; line-height: 21px; color: #333333; text-decoration: none;"></a>
												</td><td align="center">
													<a style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 14px; line-height: 21px; color: #333333; text-decoration: none;" target="_blank">&nbsp;</a>
												</td></tr>
										</table>
										<table cellpadding="0" cellspacing="0" border="0" width="100%">
											<tr>
												<td height="5"></td>
											</tr>
										</table>
										<table align="center" width="100%" border="0" cellspacing="0" cellpadding="0">
											<tr>
											</tr>
										</table>
										<div style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 14px; line-height: 21px; color: #9299a2;"> Вы получили это письмо, так как регистрируетесь на сайте ByTrip.ru</div>
										<a style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 14px; line-height: 21px; color: #9299a2; text-decoration: underline;" target="_blank">&nbsp;</a>
										<div style="font-family: -apple-system, 'Segoe UI', 'Helvetica Neue', Helvetica, Roboto, Arial, sans-serif; font-size: 14px; line-height: 21px; color: #9299a2;"> © 2023 ByTrip</div>
										<table cellpadding="0" cellspacing="0" border="0" width="100%">
											<tr>
												<td height="30"></td>
											</tr>
										</table>
									</td>
								</tr>
							</table>
							<!--[if (gte mso 9)|(IE)]>
							</td></tr></table>
							<![endif]-->
						</td>
					</tr>
					<tr em="block">
						<td align="center">
							<!--[if (gte mso 9)|(IE)]>
							<table cellpadding="0" cellspacing="0" border="0" width="660"><tr><td>
							<![endif]-->
							<table align="center" width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 660px;">
								<tr>
									<td align="center" valign="top">
										<table cellpadding="0" cellspacing="0" border="0" width="100%">
											<tr></tr>
										</table>
									</td>
								</tr>
							</table>
							<!--[if (gte mso 9)|(IE)]>
							</td></tr></table>
							<![endif]-->
						</td>
					</tr>
				</table>
				<!--[if (gte mso 9)|(IE)]>
				</td></tr></table>
				<![endif]-->
			</td>
		</tr>
	</table>
</body>
</html>
{{end}}
`

func SendMail(receiver string, authCode int) error {
	d := gomail.NewDialer(
		viper.GetString("email.host"),
		viper.GetInt("email.port"),
		viper.GetString("email.username"),
		viper.GetString("email.password"),
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	var b bytes.Buffer

	tmpl, err := template.New("template").Parse(message)
	if err != nil {
		return fmt.Errorf(fmt.Errorf("create template: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	err = tmpl.ExecuteTemplate(&b, "authCode", authCode)

	if err != nil {
		return fmt.Errorf(fmt.Errorf("execute template: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("email.username"))
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", "Подтверждение регистрации")
	m.SetBody("text/html", b.String())

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf(fmt.Errorf("dial and send: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}
