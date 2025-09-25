(self.webpackChunk_N_E = self.webpackChunk_N_E || []).push([[4999], {
    18880: e => {
        e.exports = {
            userInformation: "hVWJNc_af83e17fcda2d57128b8",
            usernameWrapper: "w9OJ1e_1ebc9dfc436405473e76",
            avatar: "jm35zc_9cc0e467b3f5bcd3e671",
            userDescriptionText: "smNw3N_88873206a1910ac23f3e",
            userDescriptionUsername: "lNyJYL_05308b9ba1c31bb20692",
            usernameCursor: "scy3zW_ad49fe0637994ca39f1b",
            userDescription: "vlbWbb_78d0d49e5bb6aca9a702",
            topRightWrapper: "mWv0Yt_a9de52ea1adc6d9927d5",
            topRightDivider: "g12LWX_62b9c1eed97eed149ffd",
            socialsWrapper: "jzObUB_485c79f8170c6a8b32ac",
            discordPresenceWidget: "vAGZGc_46e63267a2a7bcf01d88",
            discordAvatar: "w53z90_2761be3309a8ce9a65c7",
            discordStatus: "dmcmFz_023e1243733438a5e05e",
            discordStatusIcon: "mmY3CN_777ff2891e84707d8c61",
            discordUserInformation: "gc1GNb_fe2d50959b4daeee9622",
            discordActivity: "ntG92m_0efefe8489eba40080bb",
            discordUser: "eYN55m_70efc2bf41505f675a7a",
            discordUserDiv: "xGbBbN_2b1f5a69e0946a7c01de",
            discordUserBadges: "oX52tY_bf751620242550f38902",
            highlight: "z1c5YX_bd75eca7e69a5ab75a86",
            activityImage: "gzXYb3_283c0ac4075e6ceeceb8",
            discordBadge: "qLwCcm_dd68a0d61aabd8d9ab2e",
            userProfileSvg: "a1UX2N_2f8d0a9758353909192b",
            discordNotConnected: "bluuNO_93b88d0815708715c878",
            offlineText: "cw0bc2_f9265e1058509bcfff9f",
            avatarWrapper: "gcW1u9_8c5665ea139e40bf2ed2",
            discordEmoji: "kLWwVb_d6828ec770f93f1f64e9",
            discordLoading: "qYcXnV_fd673857ee90cf0d5158"
        }
    }, 21614: (e, t, s) => {
        "use strict";
        s.d(t, {A: () => i});
        var n = s(95155), a = s(97745), r = s.n(a), o = s(16990);

        function i(e) {
            let {featureName: t, onChangeFunction: s, defaultValue: a, value: i, marks: c, ...l} = e;
            return (0, n.jsxs)("div", {
                className: r().inputContainerWrapper,
                children: [(0, n.jsx)("h1", {
                    className: r().featureName,
                    children: t
                }), (0, n.jsx)("div", {
                    className: r().inputWrapperDiv,
                    children: (0, n.jsx)(o.A, {defaultValue: a, onChange: s, value: i, color: "gray", marks: c, ...l})
                })]
            })
        }
    }, 46533: e => {
        e.exports = {
            userInformation: "dGF5Z9_9c94641554f4f71d4946",
            avatar: "h03RXG_e0c9184eee4dcda6ead5",
            description: "swcRNb_651f1841b9d6b3799edd",
            username: "uVNY23_ea700dce1f1b6b7c7488",
            usernameCursor: "e30U1N_b29bd833825d560af0da",
            userDescription: "zM9yBT_41355decbc3a2d6fc806",
            badgeContainerWrapper: "cYcczm_56a083c7d3dc1184dfc3",
            socialWrapper: "pbcxwF_a92afea10069e94f042c",
            buttonWrapper: "kG0yvH_46f880e2d541e7943840",
            userButtons: "nnw01X_f9cb1c74119d032e30bb",
            button: "qYbeV1_c932c2c0eac7463998a9",
            buttonShowUrl: "iMd0Ny_dc4a04a974acbb9d5ecf",
            buttonInformationWrapper: "w35x2l_e9d17f48524ca7102ae5",
            buttonIcon: "s31Vl1_745c9bf5d5fc98cc4510",
            buttonIconWrapper: "fbXweG_8c473b71e2d8d2407e36",
            audioPlayerWrapper: "gY3w0F_144cc21dbc68909dc2c4"
        }
    }, 74999: (e, t, s) => {
        "use strict";
        s.d(t, {default: () => ev, pX: () => eb, ZI: () => eg, CO: () => ex, dz: () => ef});
        var n = s(95155), a = s(12115), r = s(75521), o = s.n(r);
        let i = e => {
            let {
                    backgroundUrl: t,
                    shuffleAudios: s,
                    audio: r,
                    videoTagRef: i,
                    audioTagRef: c,
                    audioData: l,
                    setAudioData: d
                } = e, [u, m] = (0, a.useState)(!1), p = (0, a.useRef)(null), f = "string" == typeof r,
                h = f ? "" !== r : r.length > 0;
            return (0, a.useEffect)(() => {
                (function (e) {
                    for (let t of [".mp4", ".mov", ".avi", ".m4v"]) if (e.endsWith(t)) return {video: !0};
                    return {video: !1}
                })(t).video ? m(!0) : "" !== t && (m(!1), p.current && (p.current.style.backgroundImage = "url(".concat(t, ")")))
            }, [t, p]), (0, a.useEffect)(() => {
                var e, t, n, a;
                if (f || 0 !== r.length) {
                    if (f) d({url: r, title: "Audio #1", cover: "", id: 0}); else if (1 === r.length) d({
                        url: r[0].url,
                        title: r[0].title,
                        cover: null !== (e = r[0].cover) && void 0 !== e ? e : "",
                        id: r[0].id
                    }); else if (s) {
                        let e = r[Math.floor(Math.random() * r.length)];
                        d({
                            url: e.url,
                            title: e.title,
                            cover: null !== (t = e.cover) && void 0 !== t ? t : "",
                            id: e.id
                        })
                    } else {
                        let e = !f && h && r.find(e => !0 === e.selected);
                        d({
                            url: e ? e.url : r[0].url,
                            title: e ? e.title : r[0].title,
                            cover: e ? null !== (n = e.cover) && void 0 !== n ? n : "" : null !== (a = r[0].cover) && void 0 !== a ? a : "",
                            id: e ? e.id : r[0].id
                        })
                    }
                }
            }, [r, s, f, h]), (0, n.jsxs)(n.Fragment, {
                children: [u && "" !== t ? (0, n.jsx)("video", {
                    src: t,
                    playsInline: !0,
                    loop: !0,
                    muted: h,
                    className: o().backgroundVideo,
                    ref: i
                }) : (0, n.jsx)("div", {
                    className: o().backgroundImage,
                    style: {backgroundImage: "url(".concat(t, ")")},
                    ref: p
                }), h && "" !== l.url && (0, n.jsx)("audio", {src: l.url, loop: !0, ref: c})]
            })
        };

        function c(e) {
            let {
                videoTagRef: t,
                audioTagRef: s,
                clickToEnterText: r,
                userFont: i,
                userContainerRef: c,
                animationClass: l,
                setAudioPlayer: d,
                audioContainerRef: u,
                discordWidgetRef: m
            } = e, p = (0, a.useRef)(null), f = function (e) {
                let t = arguments.length > 1 && void 0 !== arguments[1] && arguments[1],
                    s = arguments.length > 2 && void 0 !== arguments[2] ? arguments[2] : 5,
                    n = arguments.length > 3 && void 0 !== arguments[3] ? arguments[3] : 300, a = 0, r = () => {
                        e.current ? (e.current.volume = .25, e.current.play().catch(() => {
                            a < s ? (a++, setTimeout(r, n)) : console.warn("Failed to play media after ".concat(s, " attempts."))
                        }), t && (d({
                            playing: !0,
                            duration: e.current.duration,
                            currentTime: e.current.currentTime,
                            ref: e.current
                        }), e.current.addEventListener("timeupdate", () => {
                            d(t => {
                                var s, n;
                                return {
                                    ...t,
                                    currentTime: null === (s = e.current) || void 0 === s ? void 0 : s.currentTime,
                                    duration: null === (n = e.current) || void 0 === n ? void 0 : n.duration,
                                    ref: e.current
                                }
                            })
                        }))) : a < s ? (a++, setTimeout(r, n)) : console.warn("Media element not available after multiple attempts.")
                    };
                r()
            };
            return (0, n.jsx)(n.Fragment, {
                children: (0, n.jsx)("div", {
                    className: o().clickToEnterOverlay,
                    onClick: () => {
                        var e, n, a, r;
                        p.current && (p.current.classList.toggle(o().hide), p.current.style.pointerEvents = "none", setTimeout(() => {
                            p.current && (p.current.style.display = "none")
                        }, 650)), l === o().fadeUp ? (c && c.current.classList.add(l), u && (null === (e = u.current) || void 0 === e || e.classList.add(l)), m && (null === (n = m.current) || void 0 === n || n.classList.add(l))) : l === o().unfold ? document.querySelectorAll(".".concat(o().unfoldStart)).forEach((e, t) => {
                            setTimeout(() => {
                                e.classList.add(o().unfold)
                            }, 300 * t)
                        }) : l === o().scaleFade && (c && c.current.classList.add(l), u && (null === (a = u.current) || void 0 === a || a.classList.add(l)), m && (null === (r = m.current) || void 0 === r || r.classList.add(l))), f(t), f(s, !0)
                    },
                    ref: p,
                    children: (0, n.jsx)("div", {className: o().clickToEnterText, style: i.clickToEnter, children: r})
                })
            })
        }

        let l = {
                clickToEnter: {fontWeight: 700, fontSize: "30px"},
                username: {fontWeight: 600, fontSize: "36px"},
                description: {fontWeight: 570, fontSize: "17px"},
                joinDate: {fontWeight: 400, fontSize: "13px"},
                profileViewsCount: {fontWeight: 550, fontSize: "14px"},
                discordUsername: {fontWeight: 550, fontSize: "18.5px"},
                discordStatus: {fontWeight: 300, fontSize: "13.5px"},
                serverName: {fontWeight: 450, fontSize: "18px"},
                memberCount: {fontWeight: 450, fontSize: "12.5px"},
                joinButton: {fontWeight: 450, fontSize: "12px"},
                joinButtonLarge: {fontWeight: 450, fontSize: "16px"},
                buttonText: {fontWeight: 450, fontSize: "17px"},
                buttonUrl: {fontWeight: 400, fontSize: "14px"},
                customUrlHeading: {fontWeight: 450, fontSize: "18px"},
                customUrlLink: {fontWeight: 400, fontSize: "15px"},
                customUrlButtons: {fontWeight: 450, fontSize: "16px"},
                secondTabWidgetHeading: {fontWeight: 450, fontSize: "18px"},
                secondTabWidgetDescription: {fontWeight: 450, fontSize: "12px"},
                secondTabWidgetButton: {fontWeight: 450, fontSize: "11.5px"},
                audioPlayerTitle: {fontWeight: 500, fontSize: "17px"},
                audioPlayerControls: {fontWeight: 400, fontSize: "16.5px"},
                locationMarker: {fontWeight: 450, fontSize: "15px"}
            },
            d = new Set(["satoshi", "firacode", "poppins", "array", "drippy", "minecraft", "outfit", "chillax", "velocity", "technor", "jetbrains_mono", "silkscreen", "parkinsans"]),
            u = {
                satoshi: {
                    username: {fontWeight: 600, fontSize: "39.5px"},
                    description: {fontWeight: 570, fontSize: "19px"}
                },
                firacode: {
                    clickToEnter: {fontWeight: 700, fontSize: "29px"},
                    username: {fontWeight: 550, fontSize: "34px"},
                    description: {fontWeight: 570, fontSize: "15.5px"},
                    audioPlayerTitle: {fontWeight: 500, fontSize: "16px"},
                    audioPlayerControls: {fontWeight: 400, fontSize: "15px"}
                },
                array: {
                    clickToEnter: {fontWeight: 500, fontSize: "34px"},
                    username: {fontWeight: 400, fontSize: "36.5px"}
                },
                drippy: {
                    clickToEnter: {fontWeight: 200, fontSize: "31px"},
                    username: {fontWeight: 400, fontSize: "33.5px"}
                },
                velocity: {
                    clickToEnter: {fontWeight: 600, fontSize: "23px"},
                    username: {fontWeight: 600, fontSize: "29px"},
                    description: {fontWeight: 400, fontSize: "13.5px"},
                    discordUsername: {fontWeight: 600, fontSize: "16px"},
                    discordStatus: {fontWeight: 300, fontSize: "10px"},
                    audioPlayerTitle: {fontWeight: 500, fontSize: "13px", marginTop: "6px"},
                    audioPlayerControls: {fontWeight: 400, fontSize: "11px"},
                    profileViewsCount: {fontWeight: 550, fontSize: "10px"},
                    locationMarker: {fontWeight: 450, fontSize: "11px"},
                    memberCount: {fontWeight: 450, fontSize: "9.5px"},
                    serverName: {fontWeight: 450, fontSize: "16px"},
                    joinDate: {fontWeight: 400, fontSize: "10px"},
                    secondTabWidgetDescription: {fontWeight: 450, fontSize: "9px"},
                    secondTabWidgetButton: {fontWeight: 450, fontSize: "10px"},
                    customUrlHeading: {fontWeight: 450, fontSize: "15px"},
                    customUrlLink: {fontWeight: 400, fontSize: "12px"},
                    customUrlButtons: {fontWeight: 450, fontSize: "13px"}
                },
                jetbrains_mono: {
                    audioPlayerControls: {fontWeight: 400, fontSize: "14px"},
                    username: {fontWeight: 600, fontSize: "34.5px"},
                    description: {fontWeight: 570, fontSize: "16px"},
                    audioPlayerTitle: {fontWeight: 500, fontSize: "16px"},
                    discordStatus: {fontWeight: 300, fontSize: "13px"}
                },
                minecraft: {
                    username: {fontWeight: 500, fontSize: "30px"},
                    description: {fontWeight: 400, fontSize: "15px"},
                    locationMarker: {fontWeight: 450, fontSize: "14px"},
                    discordUsername: {fontWeight: 600, fontSize: "15px"},
                    discordStatus: {fontWeight: 300, fontSize: "11.5px"},
                    profileViewsCount: {fontWeight: 550, fontSize: "12.5px"},
                    clickToEnter: {fontWeight: 500, fontSize: "30px"}
                }
            }, m = {1: .7, 2: 1, 3: 1.5}, p = e => {
                let {font: t, fontSize: s} = e, n = !d.has(t) && s ? m[s] : 1;
                return Object.keys(l).reduce((e, s) => {
                    var a;
                    let r = l[s], o = null === (a = u[t]) || void 0 === a ? void 0 : a[s], i = null != o ? o : r;
                    return e[s] = {...i, fontSize: "".concat(parseFloat(i.fontSize) * n, "px")}, e
                }, {})
            };
        var f = s(63554);

        function h(e) {
            let {
                backgroundEffects: t,
                cursorEffects: s,
                cursorEffectsColor: a,
                textColor: r,
                backgroundColor: i,
                cursorUrl: c
            } = e;
            return (0, n.jsxs)(n.Fragment, {
                children: ["snowflakes" === t ? (0, n.jsx)(n.Fragment, {
                    children: (0, n.jsxs)("div", {
                        "aria-hidden": "true",
                        className: o().snowflakes,
                        children: [(0, n.jsx)("div", {
                            className: o().snowflake,
                            children: (0, n.jsx)("div", {className: o().inner, children: "❅"})
                        }), (0, n.jsx)("div", {
                            className: o().snowflake,
                            children: (0, n.jsx)("div", {className: o().inner, children: "❅"})
                        }), (0, n.jsx)("div", {
                            className: o().snowflake,
                            children: (0, n.jsx)("div", {className: o().inner, children: "❅"})
                        }), (0, n.jsx)("div", {
                            className: o().snowflake,
                            children: (0, n.jsx)("div", {className: o().inner, children: "❅"})
                        }), (0, n.jsx)("div", {
                            className: o().snowflake,
                            children: (0, n.jsx)("div", {className: o().inner, children: "❅"})
                        }), (0, n.jsx)("div", {
                            className: o().snowflake,
                            children: (0, n.jsx)("div", {className: o().inner, children: "❅"})
                        }), (0, n.jsx)("div", {
                            className: o().snowflake,
                            children: (0, n.jsx)("div", {className: o().inner, children: "❅"})
                        }), (0, n.jsx)("div", {
                            className: o().snowflake,
                            children: (0, n.jsx)("div", {className: o().inner, children: "❅"})
                        }), (0, n.jsx)("div", {
                            className: o().snowflake,
                            children: (0, n.jsx)("div", {className: o().inner, children: "❅"})
                        }), (0, n.jsx)("div", {
                            className: o().snowflake,
                            children: (0, n.jsx)("div", {className: o().inner, children: "❅"})
                        }), (0, n.jsx)("div", {
                            className: o().snowflake,
                            children: (0, n.jsx)("div", {className: o().inner, children: "❅"})
                        }), (0, n.jsx)("div", {
                            className: o().snowflake,
                            children: (0, n.jsx)("div", {className: o().inner, children: "❅"})
                        })]
                    })
                }) : "rain" === t ? (0, n.jsxs)(n.Fragment, {
                    children: [(0, n.jsx)(f.default, {
                        strategy: "beforeInteractive",
                        src: "https://assets.guns.lol/guns_storm.js",
                        id: "_gunsstorm"
                    }), (0, n.jsx)(f.default, {
                        strategy: "beforeInteractive",
                        id: "_gunsstormcode",
                        children: "\n                                snowStorm.snowColor = '".concat(r, "';\n                                snowStorm.flakesMaxActive = 80;\n                                snowStorm.useTwinkleEffect = true;\n                                snowStorm.autoStart = true;\n                                snowStorm.freezeOnBlur = false;\n                                snowStorm.excludeMobile = false;\n                            ")
                    })]
                }) : "blurred" === t ? (0, n.jsx)("div", {className: o().blurredBackground}) : "tv" === t ? (0, n.jsx)("div", {className: o().oldTV}) : "night" === t && (0, n.jsx)("div", {className: o().nightTime}), (0, n.jsx)(f.default, {
                    strategy: "beforeInteractive",
                    src: "https://assets.guns.lol/cursor-effects.js",
                    id: "_cursor"
                }), "bubbles" === s ? (0, n.jsx)(f.default, {
                    strategy: "beforeInteractive",
                    id: "_cursorbubbles",
                    children: '\n                        new cursoreffects.bubbleCursor({\n                        color: ["'.concat(a, '"]\n                        });\n                    ')
                }) : "dot" == s ? (0, n.jsx)(f.default, {
                    strategy: "beforeInteractive",
                    id: "_cursordot",
                    children: '\n                        new cursoreffects.followingDotCursor({\n                        color: ["'.concat(a, '"]\n                        });\n                    ')
                }) : "snowflakes" === s ? (0, n.jsx)(f.default, {
                    strategy: "beforeInteractive",
                    id: "_cursorsnowflakes",
                    children: '\n                        new cursoreffects.snowflakeCursor({\n                        color: ["'.concat(a, '"]\n                        });\n                    ')
                }) : "particles" === s ? (0, n.jsx)(f.default, {
                    strategy: "beforeInteractive",
                    id: "_cursorparticles",
                    children: '\n                        new cursoreffects.fairyDustCursor({\n                        colors: ["'.concat(a, '"]\n                        });\n                    ')
                }) : "ghost" === s ? (0, n.jsx)(f.default, {
                    strategy: "beforeInteractive",
                    id: "_cursorghost",
                    children: '\n                        new cursoreffects.trailingCursor({particles: 15, rate: 0.6, baseImageSrc: "'.concat(c, '"});\n                    ')
                }) : "cat" === s && (0, n.jsx)(f.default, {
                    strategy: "beforeInteractive",
                    src: "https://assets.guns.lol/oneko.js",
                    id: "_ccat"
                }), "" !== c && (0, n.jsx)("style", {children: "\n                    * {\n                      cursor: url(".concat(c, ") 16 16, auto !important;\n                    }\n                    ")}), (0, n.jsx)("style", {children: "\n                ::-moz-selection {\n                    background: ".concat(r, ";\n                    color: ").concat(i, ";\n                }\n                ::selection {\n                    background: ").concat(r, ";\n                    color: ").concat(i, ";\n                }\n                ")})]
            })
        }

        var g = s(23720), x = s(97419), b = s(90367), v = s.n(b), j = s(62179), y = s.n(j), _ = s(84203), w = s(19320);
        let N = e => {
            let {config: t, premiumConfig: s} = e, r = {
                    base: o().tooltipContainer,
                    arrow: "".concat(o().toolTipArrow, " border-2"),
                    content: "bg-transparent"
                }, i = t.user_badges, c = t.custom_badges, l = i.length > 0 && i[0] && "string" == typeof i[0],
                d = c.length > 0 && "string" == typeof c[0][0], u = {
                    staff: "Staff",
                    premium: "Premium",
                    verified: "Verified",
                    donor: "Donor",
                    og: "OG",
                    server_booster: "Server Booster",
                    winner: "Winner",
                    second: "Second",
                    third: "Third",
                    imagehost_access: "Image Host",
                    bughunter: "Bug Hunter",
                    christmas_2024: "Christmas 2024",
                    easter_2025: "Easter 2025",
                    helper: "Community Helper",
                    gifter: "Gifter",
                    the_million: "The Million"
                }, m = {
                    staff: "staff",
                    premium: "premium",
                    verified: "verified",
                    donor: "donor",
                    og: "og",
                    server_booster: "serverBooster",
                    winner: "winner",
                    second: "second",
                    third: "third",
                    imagehost_access: "imagehost",
                    bughunter: "bugHunter",
                    christmas_2024: "christmas2024",
                    easter_2025: "easter2025",
                    helper: "helper",
                    gifter: "gifter",
                    the_million: "theMillion"
                }, p = {
                    Staff: "staff",
                    Premium: "premium",
                    Verified: "verified",
                    Donor: "donor",
                    OG: "og",
                    "Server Booster": "serverBooster",
                    Winner: "winner",
                    Second: "second",
                    Third: "third",
                    "Image Host": "imageHost",
                    "Bug Hunter": "bugHunter",
                    "Christmas 2024": "christmas2024",
                    "Easter 2025": "easter2025",
                    "Community Helper": "helper",
                    Gifter: "gifter",
                    "The Million": "theMillion"
                };
            return (0, n.jsx)(n.Fragment, {
                children: (0, n.jsxs)("div", {
                    className: o().badgeContainer,
                    children: [l ? i.map(e => (0, n.jsx)(a.Fragment, {
                        children: (0, n.jsx)(w.I, {
                            content: u[e],
                            closeDelay: 150,
                            showArrow: !1,
                            offset: 7,
                            classNames: r,
                            children: (0, n.jsxs)("div", {
                                className: o().badgeContainerFlex,
                                children: [s.monochrome_badges && (0, n.jsx)("style", {children: "\n                    .".concat(o()[p[u[e]]], " {\n                        color: ").concat(s.badge_color, "!important;\n                        filter: drop-shadow(").concat(s.badge_color, " 0 0 2.5px)!important;\n                    }\n                    ")}), !t.badge_glow && (0, n.jsx)("style", {children: "\n                    .".concat(o()[p[u[e]]], " {\n                        filter: none!important;\n                    }\n                    ")}), _.A[m[e]]]
                            }, e)
                        })
                    }, e)) : i.map(e => e.enabled && (0, n.jsx)(a.Fragment, {
                        children: (0, n.jsx)(w.I, {
                            content: u[e.name],
                            closeDelay: 150,
                            showArrow: !1,
                            offset: 7,
                            classNames: r,
                            children: (0, n.jsxs)("div", {
                                className: o().badgeContainerFlex,
                                children: [s.monochrome_badges && (0, n.jsx)("style", {children: "\n                    .".concat(o()[p[u[e.name]]], " {\n                        color: ").concat(s.badge_color, "!important;\n                        filter: drop-shadow(").concat(s.badge_color, " 0 0 2.5px)!important;\n                    }\n                    ")}), !t.badge_glow && (0, n.jsx)("style", {children: "\n                    .".concat(o()[p[u[e.name]]], " {\n                        filter: none!important;\n                    }\n                    ")}), _.A[m[e.name]]]
                            }, e.name)
                        })
                    }, e.name)), 0 !== c.length && (0, n.jsx)(n.Fragment, {
                        children: d ? c.map((e, t) => (0, n.jsx)(w.I, {
                            content: e[0],
                            closeDelay: 150,
                            showArrow: !1,
                            offset: 7,
                            classNames: r,
                            children: (0, n.jsx)("div", {
                                className: o().badgeContainerFlex,
                                children: (0, n.jsx)("img", {
                                    src: e[1],
                                    alt: "Custom Badge",
                                    className: o().customBadge,
                                    draggable: !1
                                })
                            })
                        }, t)) : c.map((e, t) => e.enabled && (0, n.jsx)(w.I, {
                            content: e.name,
                            closeDelay: 150,
                            showArrow: !1,
                            offset: 7,
                            classNames: r,
                            children: (0, n.jsx)("div", {
                                className: o().badgeContainerFlex,
                                children: (0, n.jsx)("img", {
                                    src: e.icon,
                                    alt: "Custom Badge",
                                    className: o().customBadge,
                                    draggable: !1
                                })
                            })
                        }, t))
                    })]
                })
            })
        };
        var S = s(23305), C = s(6874), k = s.n(C), W = s(20417), T = s(13031), A = s(60038), D = s(95139);
        let z = (0, s(55028).default)(() => s.e(7171).then(s.bind(s, 77171)), {
            loadableGenerated: {webpack: () => [77171]},
            ssr: !1
        }), I = {
            base: o().tooltipContainer,
            arrow: "".concat(o().toolTipArrow, " border-2"),
            content: "bg-transparent"
        }, U = e => {
            let {copied: t, copy: s, content: a, children: r} = e;
            return (0, n.jsx)(w.I, {
                content: t ? "Copied" : a,
                closeDelay: 600,
                isDismissable: !0,
                showArrow: !1,
                offset: 10,
                classNames: I,
                children: (0, n.jsx)("span", {onClick: s, style: {display: "flex"}, children: r})
            })
        };

        function B(e) {
            let {socials: t, isMonochrome: s, iconColor: r, socialGlow: i, font: c, fontSize: l, username: d, ...u} = e,
                m = ["eth", "ltc", "bitcoin", "monero", "solana", "xrp"], [f, {
                    open: h,
                    close: g
                }] = (0, A.j)(!1), [x, b] = (0, a.useState)({url: "", id: ""}), [v, j] = (0, a.useState)({
                    url: "",
                    id: ""
                }), [y, _] = (0, a.useState)(""), w = p({font: c, fontSize: l}), [N, C] = (0, a.useState)([]),
                I = (e, t) => {
                    _(t), e.preventDefault(), h()
                }, B = (e, t) => {
                    let {social: s, value: n, id: a} = t, r = N.includes(a);
                    "custom_url" === s || "onlyfans" === s ? (I(e, s), "custom_url" === s ? b({url: n, id: a}) : j({
                        url: n,
                        id: a
                    })) : (ex(d, a, "click", r), C(e => [...e, a]))
                }, F = e => {
                    var t;
                    let a = S.A[e.social.replace("custom_url", "customUrl")],
                        r = null !== (t = e.mode) && void 0 !== t ? t : "link", i = D.V$[e.social], c = N.includes(e.id);
                    return (0, n.jsx)(n.Fragment, {
                        children: "link" === r ? "email" === e.social ? (0, n.jsx)(z, {
                            encoded: e.value,
                            onClick: t => {
                                ex(d, e.id, "click", c), C([...N, e.id])
                            },
                            children: a
                        }) : (0, n.jsx)(k(), {
                            href: e.value,
                            target: "_blank",
                            onClick: t => B(t, e),
                            children: "monero" === e.social && s ? S.A.moneroMonochrome : "custom_url" === e.social && e.icon ? (0, n.jsx)("img", {
                                className: o().customIcon,
                                src: e.icon,
                                alt: ""
                            }) : a
                        }) : "text" === r && (0, n.jsx)(W.i, {
                            value: "email" === e.social ? atob(e.value) : e.value,
                            children: t => {
                                let {copied: r, copy: l} = t;
                                return (0, n.jsx)(U, {
                                    copied: r,
                                    copy: l,
                                    content: "Copy ".concat(i),
                                    children: (0, n.jsx)("span", {
                                        onClick: t => {
                                            ex(d, e.id, "click", c), C([...N, e.id])
                                        },
                                        style: {display: "flex"},
                                        children: "monero" === e.social && s ? S.A.moneroMonochrome : "custom_url" === e.social && e.icon ? (0, n.jsx)("img", {
                                            className: o().customIcon,
                                            src: e.icon,
                                            alt: ""
                                        }) : a
                                    })
                                })
                            }
                        })
                    })
                }, P = e => {
                    let t = N.includes(e.id);
                    return (0, n.jsx)(W.i, {
                        value: e.value, children: a => {
                            let {copied: r, copy: o} = a;
                            return (0, n.jsx)(U, {
                                copied: r,
                                copy: o,
                                content: "Copy Address",
                                children: (0, n.jsx)("span", {
                                    onClick: s => {
                                        ex(d, e.id, "click", t), C([...N, e.id])
                                    },
                                    style: {display: "flex"},
                                    children: "monero" === e.social && s ? S.A.moneroMonochrome : S.A[e.social]
                                })
                            })
                        }
                    })
                };
            return (0, n.jsxs)(n.Fragment, {
                children: [(0, n.jsx)(T.a, {
                    opened: f,
                    onClose: g,
                    title: "custom_url" === y ? "Are you sure?" : "Content Warning",
                    centered: !0,
                    children: "custom_url" === y ? (0, n.jsxs)("div", {
                        className: o().customUrlWarning,
                        children: [(0, n.jsx)("h1", {
                            style: w.customUrlHeading,
                            children: "You are going to visit"
                        }), (0, n.jsx)("span", {
                            style: w.customUrlLink,
                            children: x.url.replace("https://", "")
                        }), (0, n.jsxs)("div", {
                            className: o().modalButtons,
                            children: [(0, n.jsx)(k(), {
                                href: x.url, target: "_blank", onClick: e => {
                                    g(), ex(d, x.id, "click", N.includes(x.id)), C([...N, x.id])
                                }, style: w.customUrlButtons, children: "Visit"
                            }), (0, n.jsx)("button", {onClick: g, style: w.customUrlButtons, children: "Cancel"})]
                        })]
                    }) : "onlyfans" === y && (0, n.jsxs)("div", {
                        className: o().customUrlWarning,
                        children: [(0, n.jsx)("h1", {
                            style: w.customUrlHeading,
                            children: "18+ Content Warning"
                        }), (0, n.jsx)("span", {
                            style: {...w.customUrlLink, wordBreak: "break-word"},
                            children: "You are about to visit an OnlyFans link. This site may contain explicit content and is intended for adults only."
                        }), (0, n.jsxs)("div", {
                            className: o().modalButtons,
                            children: [(0, n.jsx)(k(), {
                                href: v.url, target: "_blank", onClick: e => {
                                    g(), ex(d, v.id, "click", N.includes(v.id)), C([...N, v.id])
                                }, style: w.customUrlButtons, children: "Continue"
                            }), (0, n.jsx)("button", {onClick: g, style: w.customUrlButtons, children: "Cancel"})]
                        })]
                    })
                }), (0, n.jsx)("div", {
                    className: o().userSocials, ...u,
                    children: t.map((e, t) => (0, n.jsxs)("div", {
                        className: "".concat(o().userSocial, " ").concat(o()[e.social]),
                        children: [s ? s && (0, n.jsx)("style", {children: "\n                                    .".concat(o()[e.social], " svg {\n                                        color: ").concat(r, "!important;\n                                        ").concat(!1 !== i ? "filter: drop-shadow(".concat(r, " 1px 0 7px)!important;") : "filter: none!important", "\n                                    }")}) : (0, n.jsx)("style", {children: "\n                                .".concat(o()[e.social], " svg {\n                                    ").concat(!i && "filter: none!important", "\n                            }\n                            ")}), m.includes(e.social) ? P(e) : F(e)]
                    }, t))
                })]
            })
        }

        var F = s(89677);

        function P(e) {
            var t, s;
            let {
                    presenceInformation: a,
                    font: r,
                    fontSize: i,
                    description: c,
                    discordUserBadges: l,
                    discordData: d,
                    loading: u
                } = e, m = {0: "Playing", 1: "Streaming", 2: "Listening to", 3: "Watching", 5: "Competing"},
                f = a.activities, h = f && f[0], x = f && f.find(e => "Spotify" === e.name),
                b = p({font: r, fontSize: i}), v = {
                    display: (null == h ? void 0 : null === (t = h.emoji) || void 0 === t ? void 0 : t.startsWith("https://")) && (null == h ? void 0 : h.type) === 4 && !x && "flex",
                    alignItems: (null == h ? void 0 : null === (s = h.emoji) || void 0 === s ? void 0 : s.startsWith("https://")) && (null == h ? void 0 : h.type) === 4 && !x && "flex-end"
                }, j = {
                    base: o().tooltipContainer,
                    arrow: "".concat(o().toolTipArrow, " border-2"),
                    content: "bg-transparent"
                };
            return (0, n.jsx)(n.Fragment, {
                children: (0, n.jsx)("div", {
                    className: o().discordPresenceContainer,
                    style: {marginTop: "20px"},
                    children: a.success && h && !u ? (0, n.jsxs)("div", {
                        className: o().discordUserInformation,
                        children: [(0, n.jsx)("div", {
                            className: o().discordAvatar,
                            children: (0, n.jsxs)("div", {
                                className: o().discordStatus,
                                children: [(0, n.jsx)("img", {
                                    src: a.avatar,
                                    alt: "Discord Avatar"
                                }), (0, n.jsx)("img", {
                                    src: "https://assets.guns.lol/".concat(a.status, ".png"),
                                    className: o().discordStatusIcon,
                                    alt: ""
                                })]
                            })
                        }), (0, n.jsxs)("div", {
                            className: o().discordActivity, children: [(0, n.jsxs)("div", {
                                className: o().discordUser,
                                children: [(0, n.jsxs)("div", {
                                    className: o().discordUserDiv,
                                    children: [(0, n.jsx)("span", {
                                        style: b.discordUsername,
                                        children: a.username
                                    }), (0, n.jsx)("div", {
                                        className: o().discordUserBadges,
                                        children: l.map((e, t) => (0, n.jsx)(w.I, {
                                            closeDelay: 150,
                                            offset: 3,
                                            content: e,
                                            classNames: j,
                                            children: (0, n.jsx)("div", {
                                                className: o().discordBadge,
                                                children: (0, n.jsx)("img", {src: eb[e], alt: "Discord Badge"})
                                            }, t)
                                        }, t))
                                    })]
                                }), (0, n.jsx)("h3", {
                                    style: {...b.discordStatus, ...v},
                                    children: 4 !== h.type && 2 !== h.type ? (0, n.jsxs)(n.Fragment, {
                                        children: [(0, n.jsx)("span", {
                                            className: o().highlight,
                                            children: m[h.type]
                                        }), " ", h.name, (0, n.jsx)("br", {}), h.details ? h.details : h.state]
                                    }) : 2 !== h.type || x ? x ? (0, n.jsxs)(n.Fragment, {
                                        children: [(0, n.jsx)("span", {
                                            className: o().highlight,
                                            children: m[x.type]
                                        }), " ", x.details, (0, n.jsx)("br", {}), x.state && "by ".concat(x.state.replace(/;/g, ","))]
                                    }) : 4 === h.type && (0, n.jsxs)(n.Fragment, {
                                        children: [h.emoji.startsWith("https://") ? (0, n.jsx)("img", {
                                            src: h.emoji,
                                            alt: "Discord Emoji",
                                            className: o().discordEmoji
                                        }) : !h.emoji.startsWith("https://") && (0, n.jsx)("span", {
                                            className: o().defaultEmoji,
                                            children: h.emoji
                                        }), (0, n.jsx)("span", {children: h.state})]
                                    }) : (0, n.jsxs)(n.Fragment, {
                                        children: [(0, n.jsx)("span", {
                                            className: o().highlight,
                                            children: m[h.type]
                                        }), " ", h.name, (0, n.jsx)("br", {}), h.details ? h.details : h.state]
                                    })
                                })]
                            }), (0, n.jsx)("div", {
                                className: o().activityImage,
                                children: 4 !== h.type && 2 !== h.type ? h.image && "{}" !== JSON.stringify(h.image) && (0, n.jsx)(w.I, {
                                    closeDelay: 150,
                                    content: "".concat(m[h.type], " ").concat(h.name),
                                    classNames: j,
                                    children: (0, n.jsx)("img", {src: h.image, alt: ""})
                                }) : 2 !== h.type || x ? x ? (0, n.jsx)(w.I, {
                                    closeDelay: 150,
                                    content: "Listening to ".concat(x.details, " ").concat(x.state && "by ".concat(x.state.replace(/;/g, ","))),
                                    classNames: j,
                                    children: (0, n.jsx)("img", {src: x.image, alt: ""})
                                }) : "offline" === h.status && null : h.image && "{}" !== JSON.stringify(h.image) && (0, n.jsx)(w.I, {
                                    closeDelay: 150,
                                    content: "".concat(m[h.type], " ").concat(h.name),
                                    classNames: j,
                                    children: (0, n.jsx)("img", {src: h.image, alt: ""})
                                })
                            })]
                        })]
                    }) : !a.success || h || u ? u ? (0, n.jsx)("div", {
                        className: o().discordLoading,
                        children: (0, n.jsxs)("h1", {children: [g.A.discord, " Loading Discord Presence..."]})
                    }) : (0, n.jsxs)("div", {
                        className: o().discordNotConnected,
                        children: [(0, n.jsxs)("h1", {
                            style: b.discordUsername,
                            children: [g.A.userNotFound, " User Not found"]
                        }), d ? (0, n.jsxs)("h3", {
                            style: b.discordStatus,
                            children: ["Join ", (0, n.jsx)(k(), {
                                href: "https://discord.gg/guns",
                                target: "_blank",
                                children: "discord.gg/guns"
                            }), " to display your status"]
                        }) : (0, n.jsxs)("h3", {
                            style: b.discordStatus,
                            children: ["Please ", (0, n.jsx)(k(), {
                                href: "https://guns.lol/account",
                                target: "_blank",
                                children: "connect"
                            }), " your Discord account"]
                        })]
                    }) : (0, n.jsxs)("div", {
                        className: o().discordUserInformation,
                        children: [(0, n.jsx)("div", {
                            className: o().discordAvatar,
                            children: (0, n.jsxs)("div", {
                                className: o().discordStatus,
                                children: [(0, n.jsx)("img", {
                                    src: a.avatar,
                                    alt: "Discord Avatar"
                                }), (0, n.jsx)("img", {
                                    src: "https://assets.guns.lol/".concat(a.status, ".png"),
                                    className: o().discordStatusIcon,
                                    alt: ""
                                })]
                            })
                        }), (0, n.jsxs)("div", {
                            className: o().discordActivity,
                            children: [(0, n.jsxs)("div", {
                                className: o().discordUser,
                                children: [(0, n.jsxs)("div", {
                                    className: o().discordUserDiv,
                                    children: [(0, n.jsx)("span", {
                                        style: b.discordUsername,
                                        children: a.username
                                    }), (0, n.jsx)("div", {
                                        className: o().discordUserBadges,
                                        children: l.map((e, t) => (0, n.jsx)(w.I, {
                                            closeDelay: 150,
                                            offset: 3,
                                            content: e,
                                            classNames: j,
                                            children: (0, n.jsx)("div", {
                                                className: o().discordBadge,
                                                children: (0, n.jsx)("img", {src: eb[e], alt: "Discord Badge"})
                                            }, t)
                                        }, t))
                                    })]
                                }), (0, n.jsx)("h3", {
                                    style: {...b.discordStatus, ...v},
                                    children: "offline" === a.status ? (0, n.jsxs)("span", {
                                        className: o().offlineText,
                                        children: ["last seen ", a.lastSeen ? (0, F.m)(new Date(1e3 * a.lastSeen)).replace("about ", "") + " ago" : "unknown"]
                                    }) : (0, n.jsx)("span", {children: "currently doing nothing"})
                                })]
                            }), (0, n.jsx)("div", {className: o().activityImage})]
                        })]
                    })
                })
            })
        }

        function L(e) {
            let {
                publicSalt: publicSalt,
                challenge: challenge,
                timestamp: timestamp,
                nonce: nonce,
                username: c,
                premiumConfig: l,
                views: d,
                containerBorderRadius: u,
                hideCounter: m,
                isUnfold: h,
                animationStartClass: x,
                isTemplatePreview: b
            } = e, [v, j] = (0, a.useState)(!1), [turnstileRes, _] = (0, a.useState)(""), [solverResult, setSolverResult] = (0, a.useState)({
                ps: publicSalt,
                _s: "",
                n: nonce,
                responseDict: {}
            }), C = {
                base: o().tooltipContainer,
                arrow: "".concat(o().toolTipArrow, " border-2"),
                content: "bg-transparent"
            }, k = p({font: l.font, fontSize: l.font_size}), W = !!u && u > 25 && .12 * u;

            async function T(e) {
                _(e)
            }

            async function getResult(e) {
                setSolverResult({ps: publicSalt, _s: e._s, n: nonce, responseDict: e}), j(!0)
            }

            return (0, a.useEffect)(() => {
                Object.assign(window, {getToken: T, getResult: getResult});
                let e = setInterval(() => {
                    v && "" !== turnstileRes && !b && (fetch("https://guns.lol/api/analytics/record", {
                        cache: "no-store",
                        body: JSON.stringify({
                            _t: turnstileRes,
                            _s: (res => {
                                let {responseDict: t, ...rest} = res;
                                return rest
                            })(solverResult),
                            _c: challenge,
                            _ps: btoa(publicSalt),
                            _ts: timestamp,
                            __hcm: solverResult.responseDict,
                            username: decodeURIComponent(c),
                            deviceType: eg(),
                            event: "view",
                            linkId: null,
                            referrer: document.referrer
                        }),
                        method: "POST"
                    }), clearInterval(e))
                }, 200)
            }, [v, turnstileRes, solverResult]), (0, n.jsxs)(n.Fragment, {
                children: [(0, n.jsx)(f.default, {
                    strategy: "beforeInteractive",
                    id: "_code",
                    children: '\n        const workerScript = `        \n            import init, { GunsSolver } from \'https://assets.guns.lol/pkg/_gunslolpow.js\';\n            self.onmessage = async function (event) {               \n                const { ps, d, c, _n, org_ts } = event.data;\n                await init();\n                const _get_s = new GunsSolver(ps, c, parseInt(d), false, _n);\n                const ts = Math.floor(Date.now() / 1000);\n                const _tsoff = ts - Number(org_ts);\n                const tsp = ts - _tsoff;\n                const _res = await _get_s.solve_pow();    \n                if (_res !== null) {\n                    self.postMessage({ _res, ps, tsp, _n, c });\n                }\n            };\n        `;\n\n        const blob = new Blob([workerScript], { type: "application/javascript" });\n        const worker = new Worker(URL.createObjectURL(blob), { type: "module" });\n\n        worker.postMessage({\n            ps: \''.concat(publicSalt, "',\n            d: '5',\n            c: '").concat(challenge, "',\n            _n: '").concat(nonce, "',\n            org_ts: '").concat(timestamp, "'\n        });\n\n        worker.onmessage = async function (event) {\n            const { _res, ps, tsp, _n, c } = event.data;\n            if (_res !== undefined) {\n                await getResult(_res);\n            }\n        };\n    ")
                }), (0, n.jsx)(f.default, {
                    strategy: "afterInteractive",
                    src: "https://challenges.cloudflare.com/turnstile/v0/api.js",
                    id: "_turnstile"
                }), (0, n.jsx)("div", {
                    className: "cf-turnstile",
                    style: {display: "none"},
                    "data-sitekey": "0x4AAAAAAAgU7T2niLQD-TLm",
                    "data-callback": "getToken"
                }), !l.hide_views && !m && (0, n.jsx)("div", {
                    className: "".concat(o().profileViews, " ").concat(h && x),
                    style: {bottom: W ? 12 + W : 12, left: W ? 15 + W : 15},
                    children: (0, n.jsx)(w.I, {
                        color: "foreground",
                        content: "Profile Views",
                        closeDelay: 150,
                        showArrow: !1,
                        offset: 2.2,
                        classNames: C,
                        children: (0, n.jsxs)("span", {
                            style: k.profileViewsCount,
                            children: [g.A.profileViews, " ", d.toLocaleString("en-US")]
                        })
                    })
                })]
            })
        }

        function E(e) {
            var t, s, a;
            let {
                    userData: r,
                    config: i,
                    premiumConfig: c,
                    presenceInformation: l,
                    animationStartClass: d,
                    isTemplatePreview: u,
                    discordLoading: m
                } = e, f = p({font: c.font, fontSize: c.font_size}), h = {
                    base: o().tooltipContainer,
                    arrow: "".concat(o().toolTipArrow, " border-2"),
                    content: "bg-transparent"
                }, x = i.use_discord_avatar && r.discord ? l.avatar : i.avatar,
                b = "" !== i.display_name ? i.display_name : r.username, j = i.user_badges,
                _ = (null === (t = r.discord) || void 0 === t ? void 0 : t.user_badges) ? j.includes("server_booster") || j.find(e => "server_booster" === e.name) ? [...r.discord.user_badges, "Server Booster"] : r.discord.user_badges : [],
                S = null == l ? void 0 : l.avatarDecoration, C = i.custom_badges,
                k = j.length > 0 && j[0] && "string" == typeof j[0], W = C.length > 0 && "string" == typeof C[0][0],
                T = j.length > 0 && k || !k && j.some(e => !0 === e.enabled),
                A = C.length > 0 && W || !W && C.some(e => !0 === e.enabled), D = d === o().unfoldStart,
                z = null !== (s = c.typewriter_speed) && void 0 !== s ? s : 5,
                I = null !== (a = c.typewriter_delete_speed) && void 0 !== a ? a : 7,
                U = void 0 === i.discord_avatar_decoration || i.discord_avatar_decoration;
            return (0, n.jsxs)(n.Fragment, {
                children: [(0, n.jsxs)("div", {
                    className: v().userInformation,
                    children: ["" !== c.banner && (0, n.jsx)("img", {
                        src: c.banner,
                        alt: "",
                        className: v().banner
                    }), "" !== i.avatar || i.use_discord_avatar && r.discord ? (0, n.jsxs)("div", {
                        className: o().avatarWrapper,
                        children: [S && U && (0, n.jsx)("img", {
                            src: S,
                            alt: "",
                            className: o().avatarDecorationLayout1,
                            style: {top: "" !== c.banner ? "21%" : "-7%"}
                        }), (0, n.jsx)("img", {
                            src: x,
                            alt: "",
                            className: "".concat(v().avatar, " ").concat("" !== c.banner && v().adjustBanner)
                        })]
                    }) : (0, n.jsx)("div", {style: {marginTop: "" !== c.banner ? "135px" : ""}}), (0, n.jsxs)("div", {
                        className: "".concat(v().userDescription, " ").concat(D && d),
                        children: [(0, n.jsx)(w.I, {
                            color: "foreground",
                            content: "UID ".concat(r.uid.toLocaleString("en-US")),
                            closeDelay: 150,
                            showArrow: !1,
                            offset: -2,
                            classNames: h,
                            children: "typewriter" === i.username_effects ? (0, n.jsxs)(n.Fragment, {
                                children: [(0, n.jsx)("style", {children: "\n                                        .username-typewriter .Typewriter {\n                                            font-size: ".concat(f.username.fontSize, ";\n                                            font-weight: ").concat(f.username.fontWeight, ";\n                                        }\n                                    ")}), (0, n.jsx)("div", {
                                    className: "username-typewriter",
                                    children: (0, n.jsx)(y(), {
                                        options: {
                                            strings: [b],
                                            autoStart: !0,
                                            loop: !0,
                                            wrapperClassName: v().userDescriptionUsername,
                                            cursorClassName: v().usernameCursor,
                                            pauseFor: 1e3,
                                            stringSplitter: ef
                                        }
                                    })
                                })]
                            }) : (0, n.jsx)("h1", {
                                className: "".concat(v().userDescriptionUsername, " ").concat("rgb" === i.username_effects && o().rainbowName),
                                style: f.username,
                                children: b
                            })
                        }), !(!T && !A) && (0, n.jsx)("div", {
                            className: v().badgeContainerWrapper,
                            children: (0, n.jsx)(N, {config: i, premiumConfig: c})
                        }), ("" !== i.description || 0 !== c.typewriter.length) && (0, n.jsx)("h3", {
                            className: v().userDescriptionText,
                            style: f.description,
                            children: c.typewriter_enabled && 0 !== c.typewriter.length ? (0, n.jsx)(y(), {
                                options: {
                                    strings: c.typewriter,
                                    autoStart: !0,
                                    loop: !0,
                                    delay: 600 / z,
                                    deleteSpeed: 600 / I,
                                    pauseFor: 900,
                                    stringSplitter: ef
                                }
                            }) : i.description
                        }), i.location && (0, n.jsxs)("span", {
                            className: o().locationMarker,
                            style: f.locationMarker,
                            children: [g.A.locationMarker, " ", i.location]
                        })]
                    }), "enabled" === i.presence && (0, n.jsx)("div", {
                        className: "".concat(o().discordPresenceWrapper, " ").concat(D && d),
                        children: (0, n.jsx)(P, {
                            loading: m,
                            presenceInformation: l,
                            font: c.font,
                            description: i.description,
                            discordUserBadges: _,
                            discordData: r.discord,
                            fontSize: c.font_size
                        })
                    }), i.socials.length > 0 && (0, n.jsx)("div", {
                        className: "".concat(v().socialsWrapper, " ").concat(D && d),
                        children: (0, n.jsx)(B, {
                            socials: i.socials,
                            isMonochrome: i.monochrome,
                            font: c.font,
                            socialGlow: i.social_glow,
                            iconColor: i.icon_color,
                            username: r.username,
                            fontSize: c.font_size
                        })
                    })]
                }), (0, n.jsx)(L, {
                    publicSalt: r._PS,
                    challenge: r._C,
                    timestamp: r._ts,
                    nonce: r._n,
                    username: r.username,
                    premiumConfig: c,
                    views: i.page_views,
                    containerBorderRadius: c.border_radius,
                    isUnfold: D,
                    animationStartClass: d,
                    isTemplatePreview: u
                })]
            })
        }

        var R = s(97170), M = s.n(R);
        let H = (e, t) => "https://open.spotify.com/embed/".concat(e, "/").concat(t), O = e => {
            let t = e.match(/\/(?:intl-[a-zA-Z]+\/)?([a-zA-Z]+)\/([a-zA-Z0-9]+)/);
            if (t && 3 === t.length) {
                let [e, s] = [t[1], t[2]];
                return {type: e, code: s}
            }
            return null
        }, V = e => {
            let {spotifyUrl: t} = e, s = O(t);
            if (s) {
                let {type: e, code: t} = s;
                return (0, n.jsx)("iframe", {
                    src: H(e, t),
                    height: 80,
                    width: "100%",
                    frameBorder: "0",
                    style: {borderRadius: "15px"},
                    allow: "encrypted-media"
                })
            }
            return (0, n.jsx)("p", {children: "Invalid Spotify URL"})
        };

        function G(e) {
            var t, s, r, o, i, c, l, d, u, m, f, h, x;
            let {config: b, premiumConfig: v} = e, j = p({
                    font: v.font,
                    fontSize: v.font_size
                }), [y, _] = (0, a.useState)({}), [w, N] = (0, a.useState)({}), [S, C] = (0, a.useState)({}), [W, T] = (0, a.useState)({}), [A, D] = (0, a.useState)({}), [z, I] = (0, a.useState)({}), [U, B] = (0, a.useState)(!0),
                F = Object.keys(v.second_tab)[0], P = !!v.second_tab_enabled && "enabled" !== b.presence,
                L = async () => {
                    try {
                        let e = "https://discord.com/api/v9/invites/".concat(Y(v.second_tab.discord), "?with_counts=true"),
                            t = await fetch(e);
                        if (!t.ok) {
                            _({
                                error: !0,
                                approximate_member_count: "0",
                                approximate_presence_count: "0",
                                guild: {name: 404 === t.status ? "Not Found" : "Unknown Error"},
                                isVerified: !1,
                                isPartnered: !1
                            }), B(!1);
                            return
                        }
                        let s = await t.json(), n = s.guild.features.includes("VERIFIED"),
                            a = s.guild.features.includes("PARTNERED");
                        _({...s, isVerified: n, isPartnered: a})
                    } catch (e) {
                        _({
                            error: !0,
                            approximate_member_count: "0",
                            approximate_presence_count: "0",
                            guild: {name: "Unknown Error"},
                            isVerified: !1,
                            isPartnered: !1
                        })
                    } finally {
                        B(!1)
                    }
                }, E = async () => {
                    let e = X(v.second_tab.github), t = await fetch("https://api.github.com/users/".concat(e));
                    if (!t.ok) {
                        N({error: !0}), B(!1);
                        return
                    }
                    let s = await t.json();
                    N({
                        avatarUrl: s.avatar_url,
                        username: s.login,
                        followers: s.followers,
                        repositories: s.public_repos,
                        error: !1
                    }), B(!1)
                }, R = async () => {
                    let e = Z(v.second_tab.telegram), t = await fetch("https://guns.lol/api/service/telegram", {
                        method: "POST",
                        headers: {"Content-Type": "application/json"},
                        body: JSON.stringify({chatId: e})
                    }), s = await t.json();
                    if (!t.ok) {
                        C({error: !0}), B(!1);
                        return
                    }
                    C({avatar: s.avatar, name: s.name, members: s.members, error: !1}), B(!1)
                }, H = async () => {
                    let e = q(v.second_tab.youtube), t = await fetch("https://guns.lol/api/service/youtube", {
                        method: "POST",
                        headers: {"Content-Type": "application/json"},
                        body: JSON.stringify({channelId: e})
                    }), s = await t.json();
                    if (!t.ok) {
                        D({error: !0}), B(!1);
                        return
                    }
                    D({
                        avatar: s.avatar,
                        handle: s.handle,
                        subscribers: s.subscribers,
                        channelId: s.channelId,
                        error: !1
                    }), B(!1)
                }, O = async () => {
                    let e = J(v.second_tab.roblox), t = await fetch("https://guns.lol/api/service/roblox", {
                        method: "POST",
                        headers: {"Content-Type": "application/json"},
                        body: JSON.stringify({userId: e})
                    }), s = await t.json();
                    if (!t.ok) {
                        T({error: !0}), B(!1);
                        return
                    }
                    T({
                        username: s.username,
                        avatarUrl: s.avatar,
                        friends: s.friends,
                        followers: s.followers,
                        error: !1
                    }), B(!1)
                }, G = async () => {
                    let e = K(v.second_tab.lastfm), t = await fetch("https://guns.lol/api/service/lastfm", {
                        method: "POST",
                        headers: {"Content-Type": "application/json"},
                        body: JSON.stringify({username: e})
                    });
                    if (!t.ok) {
                        I({error: !0}), B(!1);
                        return
                    }
                    let s = await t.json();
                    I({
                        avatarUrl: s.avatar,
                        username: s.username,
                        artistCount: s.artist_count,
                        playcount: s.playcount,
                        error: !1
                    }), B(!1)
                }, J = e => {
                    let t = e.match(/(?:roblox\.com\/users\/)([a-zA-Z0-9-]+)/);
                    return t ? t[1] : ""
                }, Y = e => {
                    let t = e.match(/(?:discord\.gg\/|discord\.com\/invite\/|discordapp\.com\/invite\/)([a-zA-Z0-9-]+)/);
                    return t ? t[1] : ""
                }, X = e => {
                    let t = e.match(/(?:github\.com\/)([a-zA-Z0-9-]+)/);
                    return t ? t[1] : ""
                }, Z = e => {
                    let t = e.match(/(?:t\.me\/)([a-zA-Z0-9-]+)/);
                    return t ? t[1] : ""
                }, q = e => {
                    let t = e.match(/(?:youtube\.com\/(?:channel\/|user\/|c\/|@))([^/?]+)/);
                    return t ? t[1] : ""
                }, K = e => {
                    let t = e.match(/(?:last\.fm\/user\/)([a-zA-Z0-9_-]+)/);
                    return t ? t[1] : ""
                };
            (0, a.useEffect)(() => {
                v.second_tab.discord && L(), v.second_tab.github && E(), v.second_tab.telegram && R(), v.second_tab.roblox && O(), v.second_tab.youtube && H(), v.second_tab.lastfm && G()
            }, [v.second_tab.discord, v.second_tab.github, v.second_tab.telegram, v.second_tab.roblox, v.second_tab.youtube, v.second_tab.lastfm]);
            let Q = {
                discord: g.A.discord,
                spotify: g.A.spotify,
                github: g.A.github,
                telegram: g.A.telegram,
                roblox: g.A.roblox,
                youtube: g.A.youtube,
                lastfm: g.A.lastfm
            };
            return (0, n.jsx)(n.Fragment, {
                children: (0, n.jsx)("div", {
                    className: M().secondTab,
                    children: U && !v.second_tab.spotify ? (0, n.jsxs)("h1", {
                        className: M().loadingText,
                        children: [Q[F], " Loading Widget..."]
                    }) : v.second_tab.discord ? (0, n.jsxs)("div", {
                        className: M().discordServer, children: [(0, n.jsxs)("div", {
                            className: M().discordInformationWrapper,
                            children: [(0, n.jsx)("div", {
                                className: M().guildIcon,
                                children: (null == y ? void 0 : null === (t = y.guild) || void 0 === t ? void 0 : t.icon) ? (0, n.jsx)("img", {
                                    src: "https://cdn.discordapp.com/icons/".concat(y.guild.id, "/").concat(y.guild.icon, ".png"),
                                    alt: ""
                                }) : (0, n.jsx)("div", {
                                    className: M().missingIcon,
                                    children: (0, n.jsx)("h1", {children: y.error ? "!" : null === (s = y.guild) || void 0 === s ? void 0 : s.name.charAt(0)})
                                })
                            }), (0, n.jsxs)("div", {
                                className: M().guildInformation,
                                children: [(0, n.jsxs)("h1", {
                                    style: j.serverName,
                                    children: [null === (r = y.guild) || void 0 === r ? void 0 : r.name, " ", y.isVerified ? g.A.verifiedGuild : y.isPartnered ? g.A.partneredGuild : ""]
                                }), (0, n.jsxs)("h3", {
                                    style: j.memberCount,
                                    children: [(0, n.jsxs)("span", {children: [g.A.status, " ", null === (o = y.approximate_presence_count) || void 0 === o ? void 0 : o.toLocaleString("en-US"), " Online"]}), (0, n.jsxs)("span", {children: [g.A.status, " ", null === (i = y.approximate_member_count) || void 0 === i ? void 0 : i.toLocaleString("en-US"), " Members"]})]
                                }), !P && (0, n.jsx)(k(), {
                                    href: "https://discord.gg/".concat(y.error ? "" : Y(v.second_tab.discord)),
                                    target: "_blank",
                                    style: j.joinButton,
                                    children: "Join"
                                })]
                            })]
                        }), P ? (0, n.jsx)(k(), {
                            href: "https://discord.gg/".concat(y.error ? "" : Y(v.second_tab.discord)),
                            className: M().joinButtonLarge,
                            target: "_blank",
                            style: j.joinButtonLarge,
                            children: "Join"
                        }) : (0, n.jsx)("div", {className: M().widgetPlaceholder}), (0, n.jsxs)("span", {
                            className: M().platform,
                            children: [g.A.discord, " Discord"]
                        })]
                    }) : v.second_tab.spotify ? (0, n.jsx)(V, {spotifyUrl: v.second_tab.spotify}) : v.second_tab.github ? (0, n.jsxs)("div", {
                        className: M().githubWidget, children: [(0, n.jsxs)("div", {
                            className: M().githubInformationWrapper,
                            children: [(0, n.jsx)("div", {
                                className: M().githubIcon,
                                children: w.error ? (0, n.jsx)("h1", {
                                    className: M().missingIcon,
                                    children: "!"
                                }) : (0, n.jsx)("img", {src: w.avatarUrl, alt: "GitHub Avatar"})
                            }), (0, n.jsxs)("div", {
                                className: M().githubInformation,
                                children: [w.error ? (0, n.jsxs)(n.Fragment, {
                                    children: [(0, n.jsx)("h1", {
                                        style: j.secondTabWidgetHeading,
                                        children: "Unknown"
                                    }), (0, n.jsx)("h3", {
                                        style: j.secondTabWidgetDescription,
                                        children: (0, n.jsx)("span", {children: "User not found."})
                                    })]
                                }) : (0, n.jsxs)(n.Fragment, {
                                    children: [(0, n.jsx)("h1", {
                                        style: j.secondTabWidgetHeading,
                                        children: w.username
                                    }), (0, n.jsxs)("h3", {
                                        style: j.secondTabWidgetDescription,
                                        children: [(0, n.jsxs)("span", {children: [g.A.followers, " ", null === (c = w.followers) || void 0 === c ? void 0 : c.toLocaleString("en-US"), " Followers"]}), (0, n.jsxs)("span", {children: [g.A.repositories, " ", null === (l = w.repositories) || void 0 === l ? void 0 : l.toLocaleString("en-US"), " Repos"]})]
                                    })]
                                }), (0, n.jsx)(k(), {
                                    href: "https://github.com/".concat(w.username),
                                    target: "_blank",
                                    style: j.secondTabWidgetButton,
                                    children: "View Profile"
                                })]
                            })]
                        }), (0, n.jsxs)("span", {className: M().platform, children: [g.A.github, " GitHub"]})]
                    }) : v.second_tab.roblox ? (0, n.jsxs)("div", {
                        className: M().robloxWidget, children: [(0, n.jsxs)("div", {
                            className: M().robloxInformationWrapper,
                            children: [(0, n.jsx)("div", {
                                className: M().robloxIcon,
                                children: W.error ? (0, n.jsx)("h1", {
                                    className: M().missingIcon,
                                    children: "!"
                                }) : (0, n.jsx)("img", {src: W.avatarUrl, alt: "Roblox Avatar"})
                            }), (0, n.jsxs)("div", {
                                className: M().robloxInformation,
                                children: [W.error ? (0, n.jsxs)(n.Fragment, {
                                    children: [(0, n.jsx)("h1", {
                                        style: j.secondTabWidgetHeading,
                                        children: "Unknown"
                                    }), (0, n.jsx)("h3", {
                                        style: j.secondTabWidgetDescription,
                                        children: (0, n.jsx)("span", {children: "Could not fetch Roblox profile."})
                                    })]
                                }) : (0, n.jsxs)(n.Fragment, {
                                    children: [(0, n.jsx)("h1", {
                                        style: j.secondTabWidgetHeading,
                                        children: W.username
                                    }), (0, n.jsxs)("h3", {
                                        style: j.secondTabWidgetDescription,
                                        children: [(0, n.jsxs)("span", {children: [g.A.friends, " ", null === (d = W.friends) || void 0 === d ? void 0 : d.toLocaleString("en-US"), " Friends"]}), (0, n.jsxs)("span", {children: [g.A.followers, " ", null === (u = W.followers) || void 0 === u ? void 0 : u.toLocaleString("en-US"), " Followers"]})]
                                    })]
                                }), (0, n.jsx)(k(), {
                                    href: "https://roblox.com/users/".concat(J(v.second_tab.roblox), "/profile"),
                                    target: "_blank",
                                    style: j.secondTabWidgetButton,
                                    children: "View Profile"
                                })]
                            })]
                        }), (0, n.jsxs)("span", {className: M().platform, children: [g.A.roblox, " Roblox"]})]
                    }) : v.second_tab.telegram ? (0, n.jsxs)("div", {
                        className: M().telegramWidget, children: [(0, n.jsxs)("div", {
                            className: M().telegramInformationWrapper,
                            children: [(0, n.jsx)("div", {
                                className: M().telegramIcon,
                                children: S.error ? (0, n.jsx)("h1", {
                                    className: M().missingIcon,
                                    children: "!"
                                }) : S.avatar ? (0, n.jsx)("img", {
                                    src: S.avatar,
                                    alt: "Telegram Avatar"
                                }) : (0, n.jsx)("h1", {className: M().missingIcon, children: S.name && S.name[0]})
                            }), (0, n.jsxs)("div", {
                                className: M().telegramInformation,
                                children: [S.error ? (0, n.jsxs)(n.Fragment, {
                                    children: [(0, n.jsx)("h1", {
                                        style: j.secondTabWidgetHeading,
                                        children: "Unknown"
                                    }), (0, n.jsx)("h3", {
                                        style: j.secondTabWidgetDescription,
                                        children: (0, n.jsx)("span", {children: "Channel not found."})
                                    })]
                                }) : (0, n.jsxs)(n.Fragment, {
                                    children: [(0, n.jsx)("h1", {
                                        style: j.secondTabWidgetHeading,
                                        children: S.name
                                    }), (0, n.jsx)("h3", {
                                        style: j.secondTabWidgetDescription,
                                        children: (0, n.jsxs)("span", {children: [g.A.followers, " ", null === (m = S.members) || void 0 === m ? void 0 : m.toLocaleString("en-US"), " Subscribers"]})
                                    })]
                                }), (0, n.jsx)(k(), {
                                    href: "https://t.me/".concat(Z(v.second_tab.telegram)),
                                    target: "_blank",
                                    style: j.secondTabWidgetButton,
                                    children: "Join"
                                })]
                            })]
                        }), (0, n.jsxs)("span", {className: M().platform, children: [g.A.telegram, " Telegram"]})]
                    }) : v.second_tab.youtube ? (0, n.jsxs)("div", {
                        className: M().youtubeWidget,
                        children: [(0, n.jsxs)("div", {
                            className: M().youtubeInformationWrapper,
                            children: [(0, n.jsx)("div", {
                                className: M().youtubeIcon,
                                children: A.error ? (0, n.jsx)("h1", {
                                    className: M().missingIcon,
                                    children: "!"
                                }) : (0, n.jsx)("img", {src: A.avatar, alt: "YouTube Avatar"})
                            }), (0, n.jsxs)("div", {
                                className: M().youtubeInformation,
                                children: [A.error ? (0, n.jsxs)(n.Fragment, {
                                    children: [(0, n.jsx)("h1", {
                                        style: j.secondTabWidgetHeading,
                                        children: "Unknown"
                                    }), (0, n.jsx)("h3", {
                                        style: j.secondTabWidgetDescription,
                                        children: (0, n.jsx)("span", {children: "Channel not found."})
                                    })]
                                }) : (0, n.jsxs)(n.Fragment, {
                                    children: [(0, n.jsx)("h1", {
                                        style: j.secondTabWidgetHeading,
                                        children: A.handle
                                    }), (0, n.jsx)("h3", {
                                        style: j.secondTabWidgetDescription,
                                        children: (0, n.jsxs)("span", {children: [g.A.followers, " ", (x = A.subscribers) >= 1e6 ? (x / 1e6).toFixed(1) + "M" : x >= 1e3 ? (x / 1e3).toFixed(1) + "K" : x, " Subscribers"]})
                                    })]
                                }), (0, n.jsx)(k(), {
                                    href: "https://youtube.com/channel/".concat(A.channelId),
                                    target: "_blank",
                                    style: j.secondTabWidgetButton,
                                    children: "View Channel"
                                })]
                            })]
                        }), (0, n.jsxs)("span", {className: M().platform, children: [g.A.youtube, " YouTube"]})]
                    }) : v.second_tab.lastfm && (0, n.jsxs)("div", {
                        className: M().lastFMWidget, children: [(0, n.jsxs)("div", {
                            className: M().lastFMInformationWrapper,
                            children: [(0, n.jsx)("div", {
                                className: M().lastFMIcon,
                                children: z.error ? (0, n.jsx)("h1", {
                                    className: M().missingIcon,
                                    children: "!"
                                }) : (0, n.jsx)("img", {src: z.avatarUrl, alt: "Last.fm Avatar"})
                            }), (0, n.jsxs)("div", {
                                className: M().lastFMInformation,
                                children: [z.error ? (0, n.jsxs)(n.Fragment, {
                                    children: [(0, n.jsx)("h1", {
                                        style: j.secondTabWidgetHeading,
                                        children: "Unknown"
                                    }), (0, n.jsx)("h3", {
                                        style: j.secondTabWidgetDescription,
                                        children: (0, n.jsx)("span", {children: "User not found."})
                                    })]
                                }) : (0, n.jsxs)(n.Fragment, {
                                    children: [(0, n.jsx)("h1", {
                                        style: j.secondTabWidgetHeading,
                                        children: z.username
                                    }), (0, n.jsxs)("h3", {
                                        style: j.secondTabWidgetDescription,
                                        children: [(0, n.jsxs)("span", {children: [g.A.playCount, " ", null === (f = Number(z.playcount)) || void 0 === f ? void 0 : f.toLocaleString("en-US"), " Scrobbles"]}), (0, n.jsxs)("span", {children: [g.A.artistCount, " ", null === (h = Number(z.artistCount)) || void 0 === h ? void 0 : h.toLocaleString("en-US"), " Artists"]})]
                                    })]
                                }), (0, n.jsx)(k(), {
                                    href: "https://last.fm/user/".concat(K(v.second_tab.lastfm)),
                                    target: "_blank",
                                    style: j.secondTabWidgetButton,
                                    children: "View Profile"
                                })]
                            })]
                        }), (0, n.jsxs)("span", {className: M().platform, children: [g.A.lastfm, " Last.fm"]})]
                    })
                })
            })
        }

        function J(e) {
            var t, s, a;
            let {
                    userData: r,
                    config: i,
                    premiumConfig: c,
                    presenceInformation: l,
                    animationStartClass: d,
                    isTemplatePreview: u,
                    discordLoading: m
                } = e, f = p({font: c.font, fontSize: c.font_size}), h = null == l ? void 0 : l.avatarDecoration, x = {
                    base: o().tooltipContainer,
                    arrow: "".concat(o().toolTipArrow, " border-2"),
                    content: "bg-transparent"
                }, b = i.use_discord_avatar && r.discord ? l.avatar : i.avatar, v = (e => {
                    let t = new Date(1e3 * e);
                    return (0, F.m)(t)
                })(r.account_created),
                j = c.second_tab_enabled && "enabled" !== i.presence ? {display: "inherit!important"} : {},
                _ = "" !== i.display_name ? i.display_name : r.username,
                S = (null === (t = r.discord) || void 0 === t ? void 0 : t.user_badges) ? i.user_badges.includes("server_booster") || i.user_badges.find(e => "server_booster" === e.name) ? [...r.discord.user_badges, "Server Booster"] : r.discord.user_badges : [],
                C = "\n                    .".concat(M().widgetContainerWrapper, " { \n                        margin-top: 27px;\n                    } \n                    "),
                k = i.user_badges, W = i.custom_badges, T = k.length > 0 && k[0] && "string" == typeof k[0],
                A = W.length > 0 && "string" == typeof W[0][0],
                D = k.length > 0 && T || !T && k.some(e => !0 === e.enabled),
                z = W.length > 0 && A || !A && W.some(e => !0 === e.enabled),
                I = null !== (s = c.typewriter_speed) && void 0 !== s ? s : 5,
                U = null !== (a = c.typewriter_delete_speed) && void 0 !== a ? a : 7, E = i.location,
                R = d === o().unfoldStart, H = void 0 === i.discord_avatar_decoration || i.discord_avatar_decoration,
                O = c.border_radius > 35 && .15 * c.border_radius, V = Object.keys(c.second_tab).length,
                J = c.second_tab_enabled;
            return (0, n.jsxs)(n.Fragment, {
                children: [(0, n.jsxs)("div", {
                    className: "".concat(M().userProfile, " ").concat(R && d),
                    children: ["" !== i.avatar || i.use_discord_avatar && r.discord ? (0, n.jsxs)("div", {
                        className: o().avatarWrapper,
                        children: [h && H && (0, n.jsx)("img", {
                            src: h,
                            alt: "",
                            className: o().avatarDecorationLayout2
                        }), (0, n.jsx)("img", {src: b, alt: "", className: M().avatar})]
                    }) : (0, n.jsx)("span", {
                        className: M().userProfileSvg,
                        children: g.A.userAvatar
                    }), (0, n.jsxs)("div", {
                        className: M().userDescription, children: [(0, n.jsxs)("div", {
                            className: M().usernameWrapper, children: [(0, n.jsx)(w.I, {
                                color: "foreground",
                                content: "UID ".concat(r.uid.toLocaleString("en-US")),
                                closeDelay: 150,
                                showArrow: !1,
                                offset: 6,
                                classNames: x,
                                children: "typewriter" === i.username_effects ? (0, n.jsxs)(n.Fragment, {
                                    children: [(0, n.jsx)("style", {children: "\n                                        .username-typewriter .Typewriter {\n                                            font-size: ".concat(f.username.fontSize, ";\n                                            font-weight: ").concat(f.username.fontWeight, ";\n                                            line-height: 41px;\n                                        }\n                                    ")}), (0, n.jsx)("div", {
                                        className: "username-typewriter ".concat(M().usernameTypewriter),
                                        children: (0, n.jsx)(y(), {
                                            options: {
                                                strings: [_],
                                                autoStart: !0,
                                                loop: !0,
                                                wrapperClassName: M().usernameStyles,
                                                cursorClassName: M().usernameCursor,
                                                pauseFor: 1e3,
                                                stringSplitter: ef
                                            }
                                        })
                                    })]
                                }) : (0, n.jsx)("h1", {
                                    style: {...f.username, lineHeight: "41px"},
                                    className: "".concat("rgb" === i.username_effects && o().rainbowName),
                                    children: _
                                })
                            }), !(!D && !z) && (0, n.jsx)(N, {config: i, premiumConfig: c})]
                        }), ("" !== i.description || 0 !== c.typewriter.length) && (0, n.jsx)("h3", {
                            className: M().userDescriptionText,
                            style: f.description,
                            children: c.typewriter_enabled && 0 !== c.typewriter.length ? (0, n.jsx)(y(), {
                                options: {
                                    strings: c.typewriter,
                                    autoStart: !0,
                                    loop: !0,
                                    delay: 600 / I,
                                    deleteSpeed: 600 / U,
                                    pauseFor: 900,
                                    stringSplitter: ef
                                }
                            }) : i.description
                        }), !c.hide_join_date && (0, n.jsxs)("h2", {
                            style: f.joinDate,
                            children: ["joined ", v, " ago"]
                        })]
                    })]
                }), c.second_tab_enabled || "enabled" === i.presence ? c.second_tab_enabled && "enabled" === i.presence ? (0, n.jsx)("style", {
                    children: "\n                    .".concat(o().discordPresenceContainer, " {\n                        max-width: initial;\n                    }\n                    .").concat(o().highlight, " {\n                        display: none;\n                    }\n                    .").concat(o().discordActivity, " {\n                        justify-content: space-between;\n                        width: 100%;\n                    }\n                    @media (max-width: 720px) {\n                        .").concat(o().userContainer, " {\n                            width: 95%;\n                            box-sizing: border-box;\n                        }\n                        .").concat(o().audioPlayer, " {\n                            width: 95%;\n                        }\n                        .").concat(o().discordUser, " {\n                            max-width: 100%\n                        }\n                        .").concat(o().highlight, " {\n                            display: inline;\n                        }\n                        ").concat(c.second_tab_enabled && "enabled" !== i.presence ? C : "", "\n\n                        .").concat(M().secondTab, " {\n                            margin-top: 0px!important;\n                        }\n                    }\n                    @media (max-width: 430px) {\n                        .").concat(o().userContainer, " {\n                            width: 100%!important;\n                        }\n                        .").concat(o().audioPlayer, " {\n                            width: 100%!important;\n                        }\n                    }\n                    ")
                }) : (0, n.jsx)("style", {children: "\n                    .".concat(o().discordPresenceContainer, " {\n                        max-width: initial;\n                    }\n                    .").concat(o().discordActivity, " {\n                        justify-content: space-between;\n                        width: 100%;\n                    }\n                    @media (max-width: 720px) {\n                        ").concat("enabled" !== i.presence && C, "\n                    }\n                    ")}) : (0, n.jsxs)(n.Fragment, {children: [(0, n.jsx)("div", {className: M().profileSeperator}), (0, n.jsx)("style", {children: "\n                        @media (max-width: 745px) {\n                            .".concat(o().userContainer, " {\n                                width: 95%;\n                                box-sizing: border-box;\n                            }\n                        }\n                        @media (max-width: 400px) {\n                            .").concat(o().userContainer, " {\n                                width: 100%;\n                            }\n                        }\n                            \n                        ")})]}), (0, n.jsx)("div", {
                    className: "".concat(M().widgetContainerWrapper, " ").concat(R && d),
                    style: j,
                    children: (0, n.jsxs)("div", {
                        className: M().widgetContainer,
                        style: j,
                        children: ["enabled" === i.presence && (0, n.jsx)("div", {
                            className: M().widgetContainerFlex,
                            children: (0, n.jsx)(P, {
                                loading: m,
                                presenceInformation: l,
                                font: c.font,
                                description: i.description,
                                discordData: r.discord,
                                discordUserBadges: J && V >= 1 ? [] : S,
                                fontSize: c.font_size
                            })
                        }), c.second_tab_enabled && Object.keys(c.second_tab).length > 0 && (0, n.jsx)("div", {
                            className: M().widgetContainerFlex,
                            style: j,
                            children: (0, n.jsx)(G, {config: i, premiumConfig: c})
                        })]
                    })
                }), i.socials.length > 0 && (0, n.jsx)("div", {
                    className: "".concat(M().socialsWrapper, " ").concat(R && d),
                    style: {marginBottom: E ? "35px" : "20px"},
                    children: (0, n.jsx)(B, {
                        socials: i.socials,
                        isMonochrome: i.monochrome,
                        font: c.font,
                        socialGlow: i.social_glow,
                        iconColor: i.icon_color,
                        username: r.username,
                        fontSize: c.font_size
                    })
                }), (0, n.jsxs)("div", {
                    className: "".concat(M().bottomLeftWrapper, " ").concat(o().bottomLeftWrapperGlobal, " ").concat(E && o().bottomLeftWrapperLocation, " ").concat(R && d),
                    style: {bottom: O ? 12 + O : 13, left: O ? 15 + O : 18},
                    children: [(0, n.jsx)(L, {
                        publicSalt: r._PS,
                        challenge: r._C,
                        timestamp: r._ts,
                        nonce: r._n,
                        username: r.username,
                        premiumConfig: c,
                        views: i.page_views,
                        containerBorderRadius: null,
                        isUnfold: "" === E && R,
                        animationStartClass: d,
                        isTemplatePreview: u
                    }), E && (0, n.jsxs)(n.Fragment, {
                        children: [!c.hide_views && (0, n.jsx)("div", {className: M().bottomLeftDivider}), (0, n.jsxs)("span", {
                            className: o().locationMarker,
                            style: f.locationMarker,
                            children: [g.A.locationMarker, " ", E]
                        })]
                    })]
                })]
            })
        }

        var Y = s(46533), X = s.n(Y), Z = s(26696);

        function q(e) {
            let {buttons: t, showUrl: s, font: a, fontSize: r, textAlign: o} = e, i = p({font: a, fontSize: r});
            return (0, n.jsx)("div", {
                className: X().userButtons,
                children: t.map((e, t) => (0, n.jsx)("div", {
                    className: X().button,
                    style: {height: s ? 65 : ""},
                    children: (0, n.jsxs)(k(), {
                        href: e.button_url,
                        target: "_blank",
                        children: ["" !== e.button_icon && (0, n.jsx)("div", {
                            className: X().buttonIconWrapper,
                            children: (0, n.jsx)("img", {src: e.button_icon, className: X().buttonIcon, alt: ""})
                        }), (0, n.jsxs)("div", {
                            className: X().buttonInformationWrapper,
                            style: {marginLeft: "" != e.button_icon && "left" === o ? 60 : ""},
                            children: [(0, n.jsx)("h1", {
                                style: i.buttonText,
                                children: e.button_title
                            }), s && (0, n.jsxs)("span", {
                                className: X().buttonShowUrl,
                                style: i.buttonUrl,
                                children: [Z.A.link, " ", e.button_url]
                            })]
                        })]
                    })
                }, t))
            })
        }

        function K(e) {
            let {
                    animationStartClass: t,
                    audioPlayer: s,
                    setAudioData: r,
                    audioData: i,
                    setAudioPlayer: c,
                    audioContainerRef: l,
                    currentFont: d,
                    audioList: u,
                    setVolume: m,
                    fontSize: f
                } = e, h = e => {
                    if (isNaN(e) || e === 1 / 0) return "--:--";
                    let t = Math.floor(e / 60), s = Math.floor(e % 60);
                    return isNaN(t) || isNaN(s) ? "--:--" : "".concat(t, ":").concat(s < 10 ? "0".concat(s) : s)
                }, g = p({font: d, fontSize: f}), x = h(s.currentTime), b = h(s.duration),
                v = s.currentTime / s.duration * 100, j = (0, a.useRef)(!1), y = () => {
                    if (!s.ref) return;
                    let e = u.findIndex(e => e.id === i.id);
                    if (1 === u.length) {
                        s.ref.currentTime = 0;
                        return
                    }
                    w(0 === e ? u[u.length - 1] : u[e - 1])
                }, _ = () => {
                    if (!s.ref) return;
                    let e = u.findIndex(e => e.id === i.id);
                    if (1 === u.length) {
                        s.ref.currentTime = 0;
                        return
                    }
                    w(e === u.length - 1 ? u[0] : u[e + 1])
                }, w = e => {
                    let t = s.ref;
                    if (!t) return;
                    let n = () => {
                        c(e => ({...e, currentTime: t.currentTime, duration: t.duration, ref: t}))
                    }, a = () => {
                        t.play().then(() => {
                            c(e => ({...e, playing: !0, ref: t})), m({
                                muted: !1,
                                volume: 100
                            }), t.addEventListener("timeupdate", n)
                        }).catch(e => console.error("Autoplay failed:", e)), t.removeEventListener("canplaythrough", a)
                    };
                    t.pause(), t.removeEventListener("timeupdate", n), t.removeEventListener("canplaythrough", a), t.src.includes(e.url) || (t.src = e.url), t.currentTime = 0, t.volume = .25, t.loop = !0, r({
                        url: e.url,
                        title: e.title,
                        cover: e.cover,
                        id: e.id
                    }), c(e => ({
                        ...e,
                        playing: !1,
                        currentTime: 0,
                        duration: 0,
                        ref: t
                    })), t.readyState >= 3 ? a() : t.addEventListener("canplaythrough", a)
                }, N = e => {
                    j.current || (j.current = !0, e(), setTimeout(() => {
                        j.current = !1
                    }, 250))
                };
            return (0, n.jsxs)("div", {
                className: "".concat(o().audioPlayer, " ").concat(t),
                ref: l,
                children: [i.cover && "" !== i.cover ? (0, n.jsx)("img", {
                    src: i.cover,
                    style: {borderRadius: "calc(var(--containerRadius) - 10px)"},
                    alt: "",
                    className: o().audioCover
                }) : (0, n.jsx)("span", {
                    className: o().audioIcon,
                    style: {borderRadius: "calc(var(--containerRadius) - 10px)"},
                    children: Z.A.audioIcon
                }), (0, n.jsxs)("div", {
                    className: "".concat(o().audioContainer),
                    children: [(0, n.jsx)("h1", {
                        className: o().audioTitle,
                        style: g.audioPlayerTitle,
                        children: i.title
                    }), (0, n.jsxs)("div", {
                        className: o().audioControls,
                        children: [(0, n.jsxs)("div", {
                            className: o().audioProgressBarWrapper,
                            children: [(0, n.jsx)("span", {
                                style: g.audioPlayerControls,
                                children: x
                            }), (0, n.jsxs)("div", {
                                className: o().audioProgressBar,
                                children: [(0, n.jsx)("div", {className: o().audioProgress}), (0, n.jsx)("div", {
                                    className: o().audioProgressPlaying,
                                    style: {width: v + "%"}
                                })]
                            }), (0, n.jsx)("span", {style: g.audioPlayerControls, children: b})]
                        }), (0, n.jsxs)("div", {
                            className: o().controlButtons,
                            children: [(0, n.jsx)("span", {
                                className: o().sideControlButton,
                                onClick: () => N(y),
                                children: Z.A.skipBackward
                            }), (0, n.jsx)("span", {
                                className: o().mainControlButton, onClick: () => {
                                    s.playing ? (s.ref.pause(), c(e => ({
                                        ...e,
                                        playing: !1
                                    }))) : (s.ref.play(), c(e => ({...e, playing: !0})))
                                }, children: s.playing ? Z.A.pause : Z.A.play
                            }), (0, n.jsx)("span", {
                                className: o().sideControlButton,
                                onClick: () => N(_),
                                children: Z.A.skipForward
                            })]
                        })]
                    })]
                })]
            })
        }

        function Q(e) {
            var t, s;
            let {
                    userData: a,
                    config: r,
                    premiumConfig: i,
                    presenceInformation: c,
                    animationStartClass: l,
                    setAudioData: d,
                    audioData: u,
                    audioPlayer: m,
                    setAudioPlayer: f,
                    audioContainerRef: h,
                    setVolume: x
                } = e, b = p({font: i.font, fontSize: i.font_size}), v = {
                    base: o().tooltipContainer,
                    arrow: "".concat(o().toolTipArrow, " border-2"),
                    content: "bg-transparent"
                }, j = r.use_discord_avatar && a.discord ? c.avatar : r.avatar,
                _ = "" !== r.display_name ? r.display_name : a.username, S = null == c ? void 0 : c.avatarDecoration,
                C = r.user_badges, k = r.custom_badges, W = C.length > 0 && C[0] && "string" == typeof C[0],
                T = k.length > 0 && "string" == typeof k[0][0],
                A = C.length > 0 && W || !W && C.some(e => !0 === e.enabled),
                D = k.length > 0 && T || !T && k.some(e => !0 === e.enabled),
                z = null !== (t = i.typewriter_speed) && void 0 !== t ? t : 5,
                I = null !== (s = i.typewriter_delete_speed) && void 0 !== s ? s : 7,
                U = void 0 === r.discord_avatar_decoration || r.discord_avatar_decoration, F = l === o().unfoldStart;
            return (0, n.jsx)(n.Fragment, {
                children: (0, n.jsxs)("div", {
                    className: X().userInformation,
                    children: [("" !== r.avatar || r.use_discord_avatar && a.discord) && (0, n.jsxs)("div", {
                        className: o().avatarWrapper,
                        children: [S && U && (0, n.jsx)("img", {
                            src: S,
                            alt: "",
                            className: o().avatarDecorationLayout3
                        }), (0, n.jsx)("img", {src: j, alt: "", className: X().avatar})]
                    }), (0, n.jsxs)("div", {
                        className: "".concat(X().userDescription, " ").concat(F && l),
                        children: [(0, n.jsx)(w.I, {
                            color: "foreground",
                            content: "UID ".concat(a.uid.toLocaleString("en-US")),
                            closeDelay: 150,
                            showArrow: !1,
                            offset: -2,
                            classNames: v,
                            children: "typewriter" === r.username_effects ? (0, n.jsxs)(n.Fragment, {
                                children: [(0, n.jsx)("style", {children: "\n                                        .username-typewriter .Typewriter {\n                                            font-size: ".concat(b.username.fontSize, ";\n                                            font-weight: ").concat(b.username.fontWeight, ";\n                                        }\n                                    ")}), (0, n.jsx)("div", {
                                    className: "username-typewriter",
                                    children: (0, n.jsx)(y(), {
                                        options: {
                                            strings: [_],
                                            autoStart: !0,
                                            loop: !0,
                                            wrapperClassName: X().username,
                                            cursorClassName: X().usernameCursor,
                                            pauseFor: 1e3,
                                            stringSplitter: ef
                                        }
                                    })
                                })]
                            }) : (0, n.jsx)("h1", {
                                className: "".concat(X().username, " ").concat("rgb" === r.username_effects && o().rainbowName),
                                style: b.username,
                                children: _
                            })
                        }), !(!A && !D) && (0, n.jsx)("div", {
                            className: X().badgeContainerWrapper,
                            children: (0, n.jsx)(N, {config: r, premiumConfig: i})
                        }), ("" !== r.description || 0 !== i.typewriter.length) && (0, n.jsx)("h3", {
                            className: X().description,
                            style: b.description,
                            children: i.typewriter_enabled && 0 !== i.typewriter.length ? (0, n.jsx)(y(), {
                                options: {
                                    strings: i.typewriter,
                                    autoStart: !0,
                                    loop: !0,
                                    delay: 600 / z,
                                    deleteSpeed: 600 / I,
                                    pauseFor: 900,
                                    stringSplitter: ef
                                }
                            }) : r.description
                        }), r.location && (0, n.jsxs)("span", {
                            className: o().locationMarker,
                            style: b.locationMarker,
                            children: [g.A.locationMarker, " ", r.location]
                        })]
                    }), r.socials.length > 0 && (0, n.jsx)("div", {
                        className: "".concat(X().socialWrapper, " ").concat(F && l),
                        children: (0, n.jsx)(B, {
                            socials: r.socials,
                            isMonochrome: r.monochrome,
                            font: i.font,
                            socialGlow: r.social_glow,
                            iconColor: r.icon_color,
                            username: a.username,
                            fontSize: i.font_size
                        })
                    }), r.audio_player && "" !== u.url && (0, n.jsx)("div", {
                        className: X().audioPlayerWrapper,
                        children: (0, n.jsx)(K, {
                            animationStartClass: l,
                            audioPlayer: m,
                            setAudioData: d,
                            audioData: u,
                            setAudioPlayer: f,
                            audioContainerRef: h,
                            currentFont: i.font,
                            audioList: r.audio,
                            setVolume: x,
                            fontSize: i.font_size
                        })
                    }), i.buttons.length > 0 && (0, n.jsx)("div", {
                        className: "".concat(X().buttonWrapper, " ").concat(F && l),
                        children: (0, n.jsx)(q, {
                            buttons: i.buttons,
                            showUrl: i.show_url,
                            font: i.font,
                            textAlign: i.text_align,
                            fontSize: i.font_size
                        })
                    })]
                })
            })
        }

        var $ = s(32621), ee = s.n($), et = s(97745), es = s.n(et);

        function en(e) {
            let {username: t} = e;
            return (0, n.jsx)("div", {
                className: es().containerWrapper,
                children: (0, n.jsxs)("div", {
                    className: es().container,
                    children: [Z.A.error, (0, n.jsx)("h1", {children: "This user is not claimed"}), (0, n.jsx)("h3", {children: "Claim this name fast by clicking on the button below!"}), (0, n.jsxs)("div", {
                        className: es().containerButtons,
                        children: [(0, n.jsx)(k(), {
                            href: "/",
                            children: "Go Home"
                        }), (0, n.jsx)(k(), {href: "/register?claim=".concat(t), children: "Claim Now!"})]
                    })]
                })
            })
        }

        function ea() {
            return (0, n.jsx)("div", {
                className: es().containerWrapper,
                children: (0, n.jsxs)("div", {
                    className: es().container,
                    children: [Z.A.error, (0, n.jsx)("h1", {children: "Uh oh! An unknown error occurred"}), (0, n.jsxs)("h3", {
                        children: ["Please report this error on our ", (0, n.jsx)(k(), {
                            href: "https://discord.gg/guns",
                            style: {textDecoration: "none", color: "#b252cc"},
                            children: "discord server"
                        }), "."]
                    }), (0, n.jsxs)("div", {
                        className: es().containerButtons,
                        children: [(0, n.jsx)(k(), {
                            href: "/",
                            children: "Go Home"
                        }), (0, n.jsx)(k(), {
                            href: "https://discord.gg/guns",
                            target: "_blank",
                            children: "Discord Server"
                        })]
                    })]
                })
            })
        }

        function er(e) {
            let {reason: t} = e;
            return (0, n.jsx)("div", {
                className: es().containerWrapper,
                children: (0, n.jsxs)("div", {
                    className: es().container,
                    children: [Z.A.error, (0, n.jsxs)("h1", {
                        children: ["This user has been banned from ", (0, n.jsx)(k(), {
                            href: "/",
                            style: {textDecoration: "none", color: "#b252cc"},
                            children: "guns.lol"
                        })]
                    }), (0, n.jsxs)("h3", {children: ["Reason: ", t]}), (0, n.jsxs)("div", {
                        className: es().containerButtons,
                        children: [(0, n.jsx)(k(), {
                            href: "/",
                            children: "Go Home"
                        }), (0, n.jsx)(k(), {
                            href: "https://discord.gg/guns",
                            target: "_blank",
                            children: "Discord Server"
                        })]
                    })]
                })
            })
        }

        function eo() {
            return (0, n.jsx)("div", {
                className: es().containerWrapper,
                children: (0, n.jsxs)("div", {
                    className: es().container,
                    children: [Z.A.error, (0, n.jsx)("h1", {children: "This user is not verified yet"}), (0, n.jsx)("h3", {children: "If that's your account, please check your email."}), (0, n.jsxs)("div", {
                        className: es().containerButtons,
                        children: [(0, n.jsx)(k(), {
                            href: "/",
                            children: "Go Home"
                        }), (0, n.jsx)(k(), {
                            href: "https://discord.gg/guns",
                            target: "_blank",
                            children: "Discord Server"
                        })]
                    })]
                })
            })
        }

        var ei = s(35695), ec = s(15366), el = s(21614), ed = s(18880), eu = s.n(ed);

        function em(e) {
            var t, s;
            let {
                userData: r,
                config: i,
                premiumConfig: c,
                presenceInformation: l,
                animationStartClass: d,
                isTemplatePreview: u
            } = e, [m, f] = (0, a.useState)(0);
            (0, a.useEffect)(() => {
                let e = () => {
                    f(window.innerWidth)
                };
                return e(), window.addEventListener("resize", e), () => window.removeEventListener("resize", e)
            }, []);
            let h = void 0 === i.discord_avatar_decoration || i.discord_avatar_decoration,
                x = i.use_discord_avatar && r.discord ? l.avatar : i.avatar,
                b = null == l ? void 0 : l.avatarDecoration, v = p({font: c.font, fontSize: c.font_size}),
                j = i.user_badges, _ = i.custom_badges, S = j.length > 0 && j[0] && "string" == typeof j[0],
                C = _.length > 0 && "string" == typeof _[0][0],
                k = j.length > 0 && S || !S && j.some(e => !0 === e.enabled),
                W = _.length > 0 && C || !C && _.some(e => !0 === e.enabled),
                T = null !== (t = c.typewriter_speed) && void 0 !== t ? t : 5,
                A = null !== (s = c.typewriter_delete_speed) && void 0 !== s ? s : 7, D = {
                    base: o().tooltipContainer,
                    arrow: "".concat(o().toolTipArrow, " border-2"),
                    content: "bg-transparent"
                }, z = d === o().unfoldStart, I = "" !== i.display_name ? i.display_name : r.username, U = i.location,
                F = c.border_radius > 35 && .15 * c.border_radius,
                P = "" !== i.avatar || i.use_discord_avatar && r.discord,
                E = void 0 !== c.avatar_radius ? c.avatar_radius : 35;
            return (0, n.jsxs)(n.Fragment, {
                children: [(0, n.jsxs)("div", {
                    className: eu().userInformation,
                    style: {marginTop: P ? -70 : 0, marginBottom: 25 * (m <= 520 && 0 === i.socials.length)},
                    children: [P && (0, n.jsxs)("div", {
                        className: "".concat(o().avatarWrapper, " ").concat(eu().avatarWrapper),
                        children: [b && h && (0, n.jsx)("img", {
                            src: b,
                            alt: "",
                            className: o().avatarDecorationLayout4
                        }), (0, n.jsx)("img", {
                            src: x,
                            alt: "",
                            className: eu().avatar,
                            style: {borderRadius: "".concat(E, "px")}
                        })]
                    }), (0, n.jsxs)("div", {
                        className: "".concat(eu().userDescription, " ").concat(z && d), children: [(0, n.jsxs)("div", {
                            className: eu().usernameWrapper,
                            children: [(0, n.jsx)(w.I, {
                                color: "foreground",
                                content: "UID ".concat(r.uid.toLocaleString("en-US")),
                                closeDelay: 150,
                                showArrow: !1,
                                offset: -2,
                                classNames: D,
                                children: "typewriter" === i.username_effects ? (0, n.jsxs)(n.Fragment, {
                                    children: [(0, n.jsx)("style", {children: "\n                                            .username-typewriter .Typewriter {\n                                                font-size: ".concat(v.username.fontSize, ";\n                                                font-weight: ").concat(v.username.fontWeight, ";\n                                            }\n                                        ")}), (0, n.jsx)("div", {
                                        className: "username-typewriter",
                                        children: (0, n.jsx)(y(), {
                                            options: {
                                                strings: [I],
                                                autoStart: !0,
                                                loop: !0,
                                                wrapperClassName: eu().userDescriptionUsername,
                                                cursorClassName: eu().usernameCursor,
                                                pauseFor: 1e3,
                                                stringSplitter: ef
                                            }
                                        })
                                    })]
                                }) : (0, n.jsx)("h1", {
                                    className: "".concat(eu().userDescriptionUsername, " ").concat("rgb" === i.username_effects && o().rainbowName),
                                    style: v.username,
                                    children: I
                                })
                            }), !(!k && !W) && (0, n.jsx)(N, {config: i, premiumConfig: c})]
                        }), ("" !== i.description || 0 !== c.typewriter.length) && (0, n.jsx)("h3", {
                            className: eu().userDescriptionText,
                            style: v.description,
                            children: c.typewriter_enabled && 0 !== c.typewriter.length ? (0, n.jsx)(y(), {
                                options: {
                                    strings: c.typewriter,
                                    autoStart: !0,
                                    loop: !0,
                                    delay: 600 / T,
                                    deleteSpeed: 600 / A,
                                    pauseFor: 900,
                                    stringSplitter: ef
                                }
                            }) : i.description
                        })]
                    })]
                }), (0, n.jsxs)("div", {
                    className: "".concat(eu().topRightWrapper, " ").concat(o().topRightWrapperGlobal, " ").concat(U && o().topRightWrapperLocation, " ").concat(z && d),
                    style: m > 520 ? {top: F ? 15 + F : 13, right: F ? 18 + F : 18} : {
                        bottom: F ? 12 + F : 13,
                        left: F ? 15 + F : 18,
                        top: "inherit",
                        right: "inherit"
                    },
                    children: [(0, n.jsx)(L, {
                        publicSalt: r._PS,
                        challenge: r._C,
                        timestamp: r._ts,
                        nonce: r._n,
                        username: r.username,
                        premiumConfig: c,
                        views: i.page_views,
                        containerBorderRadius: null,
                        isUnfold: "" === U && z,
                        animationStartClass: d,
                        isTemplatePreview: u
                    }), U && (0, n.jsxs)(n.Fragment, {
                        children: [!c.hide_views && (0, n.jsx)("div", {className: eu().topRightDivider}), (0, n.jsxs)("span", {
                            className: o().locationMarker,
                            style: v.locationMarker,
                            children: [g.A.locationMarker, " ", U]
                        })]
                    })]
                }), i.socials.length > 0 && (0, n.jsx)("div", {
                    className: "".concat(eu().socialsWrapper, " ").concat(z && d),
                    children: (0, n.jsx)(B, {
                        socials: i.socials,
                        isMonochrome: i.monochrome,
                        font: c.font,
                        socialGlow: i.social_glow,
                        iconColor: i.icon_color,
                        username: r.username,
                        fontSize: c.font_size
                    })
                })]
            })
        }

        function ep(e) {
            var t, s;
            let {
                    presenceInformation: a,
                    font: r,
                    fontSize: i,
                    discordUserBadges: c,
                    discordData: l,
                    animationStartClass: d,
                    discordWidgetRef: u,
                    loading: m
                } = e, f = {0: "Playing", 1: "Streaming", 2: "Listening to", 3: "Watching", 5: "Competing"},
                h = a.activities, x = h && h[0], b = h && h.find(e => "Spotify" === e.name),
                v = p({font: r, fontSize: i}), j = {
                    display: (null == x ? void 0 : null === (t = x.emoji) || void 0 === t ? void 0 : t.startsWith("https://")) && (null == x ? void 0 : x.type) === 4 && !b && "flex",
                    alignItems: (null == x ? void 0 : null === (s = x.emoji) || void 0 === s ? void 0 : s.startsWith("https://")) && (null == x ? void 0 : x.type) === 4 && !b && "flex-end"
                }, y = {
                    base: o().tooltipContainer,
                    arrow: "".concat(o().toolTipArrow, " border-2"),
                    content: "bg-transparent"
                };
            return (0, n.jsx)(n.Fragment, {
                children: (0, n.jsx)("div", {
                    className: "".concat(eu().discordPresenceWidget, " ").concat(d),
                    ref: u,
                    children: a.success && x && !m ? (0, n.jsxs)("div", {
                        className: eu().discordUserInformation,
                        children: [(0, n.jsx)("div", {
                            className: eu().discordAvatar,
                            children: (0, n.jsxs)("div", {
                                className: eu().discordStatus,
                                children: [(0, n.jsx)("img", {
                                    src: a.avatar,
                                    alt: "Discord Avatar"
                                }), (0, n.jsx)("img", {
                                    src: "https://assets.guns.lol/".concat(a.status, ".png"),
                                    className: eu().discordStatusIcon,
                                    alt: ""
                                })]
                            })
                        }), (0, n.jsxs)("div", {
                            className: eu().discordActivity, children: [(0, n.jsxs)("div", {
                                className: eu().discordUser,
                                children: [(0, n.jsxs)("div", {
                                    className: eu().discordUserDiv,
                                    children: [(0, n.jsx)("span", {
                                        style: v.discordUsername,
                                        children: a.username
                                    }), (0, n.jsx)("div", {
                                        className: eu().discordUserBadges,
                                        children: c.map((e, t) => (0, n.jsx)(w.I, {
                                            closeDelay: 150,
                                            offset: 3,
                                            content: e,
                                            classNames: y,
                                            children: (0, n.jsx)("div", {
                                                className: eu().discordBadge,
                                                children: (0, n.jsx)("img", {src: eb[e], alt: "Discord Badge"})
                                            }, t)
                                        }, t))
                                    })]
                                }), (0, n.jsx)("h3", {
                                    style: {...v.discordStatus, ...j},
                                    children: 4 !== x.type && 2 !== x.type ? (0, n.jsxs)(n.Fragment, {
                                        children: [(0, n.jsx)("span", {
                                            className: eu().highlight,
                                            children: f[x.type]
                                        }), " ", x.name, (0, n.jsx)("br", {}), x.details ? x.details : x.state]
                                    }) : 2 !== x.type || b ? b ? (0, n.jsxs)(n.Fragment, {
                                        children: [(0, n.jsx)("span", {
                                            className: eu().highlight,
                                            children: f[b.type]
                                        }), " ", b.details, (0, n.jsx)("br", {}), b.state && "by ".concat(b.state.replace(/;/g, ","))]
                                    }) : 4 === x.type && (0, n.jsxs)(n.Fragment, {
                                        children: [x.emoji.startsWith("https://") ? (0, n.jsx)("img", {
                                            src: x.emoji,
                                            alt: "Discord Emoji",
                                            className: eu().discordEmoji
                                        }) : !x.emoji.startsWith("https://") && (0, n.jsx)("span", {
                                            className: eu().defaultEmoji,
                                            children: x.emoji
                                        }), (0, n.jsx)("span", {children: x.state})]
                                    }) : (0, n.jsxs)(n.Fragment, {
                                        children: [(0, n.jsx)("span", {
                                            className: eu().highlight,
                                            children: f[x.type]
                                        }), " ", x.name, (0, n.jsx)("br", {}), x.details ? x.details : x.state]
                                    })
                                })]
                            }), (0, n.jsx)("div", {
                                className: eu().activityImage,
                                children: 4 !== x.type && 2 !== x.type ? x.image && "{}" !== JSON.stringify(x.image) && (0, n.jsx)(w.I, {
                                    closeDelay: 150,
                                    content: "".concat(f[x.type], " ").concat(x.name),
                                    classNames: y,
                                    children: (0, n.jsx)("img", {src: x.image, alt: ""})
                                }) : 2 !== x.type || b ? b ? (0, n.jsx)(w.I, {
                                    closeDelay: 150,
                                    content: "Listening to ".concat(b.details, " ").concat(b.state && "by ".concat(b.state.replace(/;/g, ","))),
                                    classNames: y,
                                    children: (0, n.jsx)("img", {src: b.image, alt: ""})
                                }) : "offline" === x.status && null : x.image && "{}" !== JSON.stringify(x.image) && (0, n.jsx)(w.I, {
                                    closeDelay: 150,
                                    content: "".concat(f[x.type], " ").concat(x.name),
                                    classNames: y,
                                    children: (0, n.jsx)("img", {src: x.image, alt: ""})
                                })
                            })]
                        })]
                    }) : !a.success || x || m ? m ? (0, n.jsx)("div", {
                        className: eu().discordLoading,
                        children: (0, n.jsxs)("h1", {children: [g.A.discord, " Loading Discord Presence..."]})
                    }) : (0, n.jsxs)("div", {
                        className: eu().discordNotConnected,
                        children: [(0, n.jsxs)("h1", {
                            style: v.discordUsername,
                            children: [g.A.userNotFound, " User Not found"]
                        }), l ? (0, n.jsxs)("h3", {
                            style: v.discordStatus,
                            children: ["Join ", (0, n.jsx)(k(), {
                                href: "https://discord.gg/guns",
                                target: "_blank",
                                children: "discord.gg/guns"
                            }), " to display your Discord presence"]
                        }) : (0, n.jsxs)("h3", {
                            style: v.discordStatus,
                            children: ["Please ", (0, n.jsx)(k(), {
                                href: "https://guns.lol/account",
                                target: "_blank",
                                children: "connect"
                            }), " your Discord account to display Discord presence"]
                        })]
                    }) : (0, n.jsxs)("div", {
                        className: eu().discordUserInformation,
                        children: [(0, n.jsx)("div", {
                            className: eu().discordAvatar,
                            children: (0, n.jsxs)("div", {
                                className: eu().discordStatus,
                                children: [(0, n.jsx)("img", {
                                    src: a.avatar,
                                    alt: "Discord Avatar"
                                }), (0, n.jsx)("img", {
                                    src: "https://assets.guns.lol/".concat(a.status, ".png"),
                                    className: eu().discordStatusIcon,
                                    alt: ""
                                })]
                            })
                        }), (0, n.jsxs)("div", {
                            className: eu().discordActivity,
                            children: [(0, n.jsxs)("div", {
                                className: eu().discordUser,
                                children: [(0, n.jsxs)("div", {
                                    className: eu().discordUserDiv,
                                    children: [(0, n.jsx)("span", {
                                        style: v.discordUsername,
                                        children: a.username
                                    }), (0, n.jsx)("div", {
                                        className: eu().discordUserBadges,
                                        children: c.map((e, t) => (0, n.jsx)(w.I, {
                                            closeDelay: 150,
                                            offset: 3,
                                            content: e,
                                            classNames: y,
                                            children: (0, n.jsx)("div", {
                                                className: eu().discordBadge,
                                                children: (0, n.jsx)("img", {src: eb[e], alt: "Discord Badge"})
                                            }, t)
                                        }, t))
                                    })]
                                }), (0, n.jsx)("h3", {
                                    style: {...v.discordStatus, ...j},
                                    children: "offline" === a.status ? (0, n.jsxs)("span", {
                                        className: eu().offlineText,
                                        children: ["last seen ", a.lastSeen ? (0, F.m)(new Date(1e3 * a.lastSeen)).replace("about ", "") + " ago" : "unknown"]
                                    }) : (0, n.jsx)("span", {children: "currently doing nothing"})
                                })]
                            }), (0, n.jsx)("div", {className: eu().activityImage})]
                        })]
                    })
                })
            })
        }

        let ef = e => new (ee())().splitGraphemes(e), eh = (e, t) => {
            let s = parseInt(e = e.replace(/^#/, ""), 16), n = Math.max(0, Math.min(1, t));
            return "rgba(".concat(s >> 16 & 255, ", ").concat(s >> 8 & 255, ", ").concat(255 & s, ", ").concat(n, ")")
        };

        function eg() {
            let e = navigator.userAgent.toLowerCase();
            return /mobile|iphone|ipod|android.*mobile/.test(e) ? "mobile" : /ipad|tablet|android(?!.*mobile)/.test(e) ? "tablet" : "desktop"
        }

        async function ex(e, t, s, n) {
            if (n) return;
            let a = await fetch("".concat("https://guns.lol", "/api/analytics/record"), {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({username: e, event: s, linkId: t, referrer: document.referrer, deviceType: eg()})
            }), r = await a.json();
            r.error && console.error("Error recording analytics:", r.error)
        }

        let eb = {
            "Discord Staff": "https://cdn.discordapp.com/badge-icons/5e74e9b61934fc1f67c65515d1f7e60d.png",
            "Partnered Server Owner": "https://cdn.discordapp.com/badge-icons/3f9748e53446a137a052f3454e2de41e.png",
            "Moderator Programs Alumni": "https://cdn.discordapp.com/badge-icons/fee1624003e2fee35cb398e125dc479b.png",
            "HypeSquad Events Member": "https://cdn.discordapp.com/badge-icons/bf01d1073931f921909045f3a39fd264.png",
            "HypeSquad Bravery": "https://cdn.discordapp.com/badge-icons/8a88d63823d8a71cd5e390baa45efa02.png",
            "HypeSquad Brilliance": "https://cdn.discordapp.com/badge-icons/011940fd013da3f7fb926e4a1cd2e618.png",
            "HypeSquad Balance": "https://cdn.discordapp.com/badge-icons/3aa41de486fa12454c3761e8e223442e.png",
            "Bug Hunter Level 1": "https://cdn.discordapp.com/badge-icons/2717692c7dca7289b35297368a940dd0.png",
            "Bug Hunter Level 2": "https://cdn.discordapp.com/badge-icons/848f79194d4be5ff5f81505cbd0ce1e6.png",
            "Active Developer": "https://cdn.discordapp.com/badge-icons/6bdc42827a38498929a4920da12695d9.png",
            "Early Verified Bot Developer": "https://cdn.discordapp.com/badge-icons/6df5892e0f35b051f8b61eace34f4967.png",
            "Early Nitro Supporter": "https://cdn.discordapp.com/badge-icons/7060786766c9c840eb3019e725d2b358.png",
            "Early Supporter": "https://cdn.discordapp.com/badge-icons/7060786766c9c840eb3019e725d2b358.png",
            "Discord Nitro": "https://cdn.discordapp.com/badge-icons/2ba85e8026a8614b640c2837bcdfe21b.png",
            "Server Booster": "https://cdn.discordapp.com/badge-icons/72bed924410c304dbe3d00a6e593ff59.png"
        };

        function ev(e) {
            var t, s, r, l, d, u;
            let {data: m} = e, f = (0, ei.useRouter)(),
                b = (0, ei.useSearchParams)(), [v, j] = (0, a.useState)(!1), [y, _] = (0, a.useState)({
                    templatePreview: !1,
                    templateId: "",
                    error: null
                }), [w, N] = (0, a.useState)({}), [S, C] = (0, a.useState)(!1),
                k = decodeURIComponent((0, ei.useParams)().username), [W, T] = (0, a.useState)(!1), [A, D] = (0, a.useState)({
                    volume: 100,
                    muted: !1
                }),
                z = null === (r = m.config) || void 0 === r ? void 0 : null === (s = r.premium) || void 0 === s ? void 0 : null === (t = s.metadata) || void 0 === t ? void 0 : t.favicon;
            m.alias && m.redirect && (0, ei.redirect)(m.redirect);
            let [I, U] = (0, a.useState)(null !== (u = m.config) && void 0 !== u ? u : {
                    socials: [],
                    user_badges: [],
                    custom_badges: [],
                    color: "#1b1b1b",
                    text_color: "#FFFFFF",
                    bg_color: "#080808",
                    gradient_1: "#881c9c",
                    gradient_2: "#0d0d0d",
                    profile_gradient: !1,
                    use_discord_avatar: !1,
                    volume_control: !1,
                    badge_glow: !0,
                    username_glow: !0,
                    swap_colors: !1,
                    monochrome: !0,
                    animated_title: !1,
                    page_views: 0,
                    blur: 7,
                    opacity: 1,
                    url: "",
                    audio: "",
                    avatar: "",
                    custom_cursor: "",
                    description: "",
                    presence: "",
                    display_name: "",
                    background_effects: "",
                    username_effects: "",
                    premium: {
                        effects_color: "#ffffff",
                        badge_color: "#ffffff",
                        border_color: "#ffffff",
                        typewriter_enabled: !1,
                        typewriter: [],
                        hide_views: !1,
                        border_enabled: !1,
                        second_tab_enabled: !1,
                        monochrome_badges: !1,
                        border_width: 2,
                        border_radius: 18,
                        second_tab: {},
                        cursor_effects: "",
                        font: "satoshi",
                        page_enter_text: "",
                        layout: "default",
                        banner: "",
                        buttons: [],
                        show_url: !1,
                        text_align: "center",
                        button_shadow: !1,
                        button_border_radius: 12
                    }
                }), [B, F] = (0, a.useState)(I.premium), [P, R] = (0, a.useState)({
                    playing: !1,
                    currentTime: 0,
                    duration: 0,
                    ref: null
                }), [M, H] = (0, a.useState)({url: "", title: "", cover: "", id: ""}), O = (0, a.useRef)(null),
                V = (0, a.useRef)(null), G = (0, a.useRef)(null), Y = (0, a.useRef)(null), X = (0, a.useRef)(null),
                Z = p({font: B.font, fontSize: B.font_size}), q = I.user_badges,
                $ = (null === (l = m.discord) || void 0 === l ? void 0 : l.user_badges) ? q.includes("server_booster") || q.find(e => "server_booster" === e.name) ? [...m.discord.user_badges, "Server Booster"] : m.discord.user_badges : [],
                ee = I.text_color,
                et = "".concat(B.border_width, "px solid ").concat(eh(B.border_color, I.opacity - .2 * I.opacity)),
                es = "2px solid ".concat(eh(B.border_color, I.opacity - .1 * I.opacity)),
                ed = "simplistic" === B.layout ? {justifyContent: "initial"} : void 0, eu = {
                    "--presenceUsernameColor": ee,
                    "--presenceStatusColor": eh(ee, .7),
                    "--presenceContainerBackground": I.swap_colors ? eh(I.color, .2) : eh(ee, .07),
                    "--presenceContainerBorder": I.swap_colors ? "2px solid ".concat(eh(I.color, .1)) : "2px solid ".concat(eh(ee, .05)),
                    "--discordServerMembers": eh(ee, .6),
                    "--secondTabBackground": I.swap_colors ? eh(I.color, .2) : eh(ee, .07),
                    "--secondTabBorder": I.swap_colors ? "2px solid ".concat(eh(I.color, .1)) : "2px solid ".concat(eh(ee, .05)),
                    "--discordMissingIcon": eh(ee, .25),
                    "--modernLayoutSeperator": eh(ee, .12),
                    "--modernLayoutSeperatorBorder": eh(ee, .12),
                    "--discordAvatarBorder": "2px solid ".concat(eh(I.color, .15)),
                    "--badgeContainerBackground": I.swap_colors ? eh(I.color, .2) : eh(ee, .08),
                    "--badgeContainerBorder": "2px solid ".concat(eh(I.swap_colors ? I.color : ee, .04)),
                    "--userBadge": {default: "20.5px", modern: "18px", simplistic: "20.5px", sleek: "20px"}[B.layout],
                    "--textColor": ee,
                    "--textColorDarker": eh(ee, .5),
                    "--iconColor": I.icon_color,
                    "--backgroundColor": I.bg_color,
                    "--joinStatus": eh(ee, .65),
                    "--audioPlayerBar": eh(ee, .35),
                    "--avatarBorder": B.border_enabled ? es : "2px solid ".concat(eh(I.color, .4)),
                    "--containerColor": eh(I.color, I.opacity),
                    "--containerBlur": "".concat(I.blur, "px"),
                    "--containerGradient": I.profile_gradient && "linear-gradient(25deg, ".concat(eh(I.gradient_1, I.opacity), ", ").concat(eh(I.gradient_2, I.opacity), ")"),
                    "--containerBorder": B.border_enabled && et,
                    "--containerRadius": "".concat(B.border_radius, "px"),
                    "--containerWidth": {
                        default: "40rem",
                        modern: "44rem",
                        simplistic: "35.5rem",
                        sleek: "40rem"
                    }[B.layout],
                    "--containerPadding": {
                        default: "35px 15px",
                        modern: "27px 23px",
                        simplistic: "35px 15px",
                        sleek: "25px"
                    }[B.layout],
                    "--volumeBackgroundcolor": eh(I.color, 0 !== I.opacity ? I.opacity : .2),
                    "--volumeBorder": "2px solid ".concat(eh(I.color, 0 !== I.opacity ? I.opacity - .1 * I.opacity : .05)),
                    "--colorUsernameGlow": I.username_glow && "0px 0px 16.5px ".concat(ee),
                    "--usernameEffects": !["rgb", "none", "typewriter"].includes(I.username_effects) && "url(https://assets.guns.lol/sparkle_".concat(I.username_effects, ".gif)"),
                    "--buttonBackground": eh(ee, .22),
                    "--buttonBorder": "2px solid ".concat(eh(ee, .12)),
                    "--buttonRadius": "".concat(B.button_border_radius, "px"),
                    "--buttonAlign": B.text_align,
                    "--buttonBackgroundHover": eh(ee, .32),
                    "--buttonBorderHover": "2px solid ".concat(eh(ee, .25)),
                    "--buttonUrlColor": eh(ee, .65),
                    "--buttonShadow": B.button_shadow && "".concat(eh(ee, .22), " 0px 0px 15px"),
                    "--profileViewsContainerBorder": "2px solid ".concat(eh(B.border_color, I.opacity - .2 * I.opacity)),
                    "--templatePreviewContainerBackground": eh(I.color, .4),
                    "--platformTextColor": eh(ee, .4),
                    "--secondTabWidgetDescription": eh(ee, .6),
                    "--audioIconBackground": eh(ee, .2),
                    "--progessBarFull": eh(ee, .3),
                    "--progessBarPlaying": ee,
                    "--controlButtonsOther": eh(ee, .5),
                    "--controlButtonPlaying": ee,
                    "--controlTextColor": eh(ee, .7),
                    "--modernLayoutBottomLeftDivider": eh(ee, .14),
                    "--sleekLayoutTopRightDivider": eh(ee, .14)
                },
                ef = m.premium ? "unfold" === B.animation ? o().unfoldStart : "fade" === B.animation ? o().fadeUpStart : "pop" === B.animation ? o().scaleFadeStart : o().fadeUpStart : o().fadeUpStart,
                eg = m.premium ? "unfold" === B.animation ? o().unfold : "fade" === B.animation ? o().fadeUp : "pop" === B.animation ? o().scaleFade : o().fadeUp : o().fadeUp,
                ex = {
                    modern: "".concat(o().userContainer, " ").concat(ef),
                    default: "".concat(o().userContainer, " ").concat(ef),
                    simplistic: "".concat(o().userContainerLayout3, " ").concat(ef),
                    sleek: "".concat(o().userContainerLayout4, " ").concat(ef)
                };
            (0, a.useEffect)(() => {
                if (j("true" === b.get("custom_domain")), m.success && !m.ban && I.animated_title) {
                    var e;
                    let t, s;
                    e = "" !== I.display_name ? I.display_name : m.username, t = 0, s = !0, setInterval(() => {
                        t === e.length && (s = !1), 0 === t && (s = !0);
                        let n = 0 === (t = s ? t + 1 : t - 1) ? "" : e.slice(0, t);
                        document.title = "@" + n
                    }, 380)
                }
                ev || e_(I.url) || (eg === o().fadeUp ? setTimeout(function () {
                    var e, t;
                    G.current && G.current.classList.add(o().fadeUp), Y.current && (null === (e = Y.current) || void 0 === e || e.classList.add(o().fadeUp)), X.current && (null === (t = X.current) || void 0 === t || t.classList.add(o().fadeUp))
                }, 350) : eg === o().unfold ? document.querySelectorAll(".".concat(o().unfoldStart)).forEach((e, t) => {
                    setTimeout(() => {
                        e.classList.add(o().unfold)
                    }, 300 * t)
                }) : eg === o().scaleFade && setTimeout(function () {
                    var e, t;
                    G.current && G.current.classList.add(o().scaleFade), Y.current && (null === (e = Y.current) || void 0 === e || e.classList.add(o().scaleFade)), X.current && (null === (t = X.current) || void 0 === t || t.classList.add(o().scaleFade))
                }, 350))
            }, [b, I]), (0, a.useEffect)(() => {
                let e = navigator.userAgent.toLowerCase();
                T(/iphone|ipod|android.*mobile|windows phone|blackberry|bb10/.test(e) || /ipad|tablet|android(?!.*mobile)/.test(e) || navigator.maxTouchPoints > 1 && /macintosh/.test(e))
            }, []);
            let eb = "string" == typeof I.audio, [ev, ej] = (0, a.useState)(eb ? "" !== I.audio : I.audio.length > 0),
                ey = () => {
                    let e = !A.muted, t = P.ref;
                    0 === A.volume && A.muted ? (D({
                        ...A,
                        muted: !1,
                        volume: 100
                    }), t && (t.muted = !1, t.volume = .25), !ev && O.current && (O.current.muted = !1, O.current.volume = .25)) : (D({
                        ...A,
                        muted: e
                    }), t && (t.muted = e, e || (t.volume = A.volume / 100 * .25)), ev || !O.current || (O.current.muted = e, e || (O.current.volume = A.volume / 100 * .25)))
                }, e_ = e => {
                    for (let t of [".mp4", ".mov", ".avi", ".m4v"]) if (e.endsWith(t)) return !0;
                    return !1
                };
            return (0, a.useEffect)(() => {
                let e = b.get("templatePreview"), t = b.get("templateId");
                (e && !t || !e && t) && f.push("/".concat(k));
                let s = async () => {
                    let s = await fetch("".concat("https://guns.lol", "/api/templates/preview/").concat(t), {method: "POST"}),
                        a = await s.json();
                    _({templatePreview: "true" === e, templateId: String(t), error: a.error});
                    let r = a.config, o = a.premium, {mergedConfig: i, mergedPremiumConfig: c} = n(r, o);
                    U(i), m.premium && o && F(c), ej("" !== r.audio || r.audio.length > 0)
                };
                e && t && s();
                let n = (e, t) => ({mergedConfig: {...I, ...e}, mergedPremiumConfig: {...B, ...t}})
            }, [b]), (0, a.useEffect)(() => {
                let e = async e => {
                    let t = await fetch("".concat("https://guns.lol", "/api/service/discord"), {
                        method: "POST",
                        headers: {"Content-Type": "application/json"},
                        body: JSON.stringify({type: e, discordId: m.discord.id})
                    }), s = await t.json();
                    t.ok ? N(s) : console.error("Error fetching Discord ".concat(e, ":"), s.error), C(!1)
                };
                m.discord && "enabled" === m.config.presence && !m.ban && m.verified ? (C(!0), e("presence")) : m.discord && m.config.use_discord_avatar && !m.ban && m.verified && e("avatar")
            }, [m.discord_id]), (0, n.jsxs)(n.Fragment, {
                children: [z && "" !== z && (0, n.jsx)("link", {
                    rel: "icon",
                    href: z,
                    sizes: "any"
                }), m.success && !m.ban && m.verified ? (0, n.jsxs)(n.Fragment, {
                    children: [(0, n.jsx)(i, {
                        backgroundUrl: I.url,
                        audio: I.audio,
                        shuffleAudios: I.shuffle_audios,
                        videoTagRef: O,
                        audioTagRef: V,
                        audioData: M,
                        setAudioData: H
                    }), (ev || e_(I.url)) && (0, n.jsx)(c, {
                        videoTagRef: O,
                        audioTagRef: V,
                        clickToEnterText: "" === B.page_enter_text ? "click to enter..." : B.page_enter_text,
                        userFont: Z,
                        userContainerRef: G,
                        animationClass: eg,
                        setAudioPlayer: R,
                        audioContainerRef: Y,
                        discordWidgetRef: X
                    }), (0, n.jsx)(h, {
                        backgroundEffects: I.background_effects,
                        cursorEffects: B.cursor_effects,
                        cursorEffectsColor: B.effects_color,
                        textColor: ee,
                        backgroundColor: I.bg_color,
                        cursorUrl: I.custom_cursor
                    }), (0, n.jsx)(x.A, {
                        font: B.font,
                        letterSpacing: B.font_letter_spacing
                    }), (0, n.jsx)("div", {
                        className: o().backgroundColor,
                        style: {backgroundColor: I.bg_color}
                    }), (0, n.jsxs)("div", {
                        style: eu, children: [I.volume_control && (e_(I.url) || ev) && (0, n.jsxs)("div", {
                            className: o().volumeButton,
                            children: [A.muted ? (0, n.jsx)("span", {
                                onClick: ey,
                                children: g.A.muted
                            }) : (0, n.jsx)("span", {onClick: ey, children: g.A.unmuted}), (0, n.jsxs)("div", {
                                className: o().volumeSlider,
                                children: [(0, n.jsx)("style", {children: "\n                                        .mantine-Slider-label {\n                                        display: none\n                                        }\n                                        .mantine-Slider-track::before {\n                                                background: var(--textColorDarker);\n                                            }\n                                        .mantine-Slider-bar {\n                                            background-color: var(--textColor)\n                                        }\n                                        .mantine-Slider-thumb {\n                                            border: 0.15rem solid var(--containerColor);\n                                            box-shadow: var(--textColorDarker) 0px 5px 25px 0px;\n                                        "}), (0, n.jsx)(el.A, {
                                    value: A.muted ? 0 : Math.round(A.volume),
                                    onChangeFunction: e => {
                                        let t = 0 === e, s = P.ref;
                                        D({
                                            volume: e,
                                            muted: t
                                        }), s && (s.volume = e / 100 * .25, s.muted = t), !ev && O.current && (O.current.volume = e / 100 * .25, O.current.muted = t)
                                    },
                                    step: 1,
                                    max: 100,
                                    min: 0
                                })]
                            })]
                        }), y.templatePreview && (0, n.jsx)("div", {
                            className: o().templatePreviewContainer,
                            style: {
                                backgroundColor: "var(--templatePreviewContainerBackground)",
                                backdropFilter: "blur(var(--containerBlur))"
                            },
                            children: y.error ? y.error.includes("access token") ? (0, n.jsx)("h1", {
                                style: {color: ee},
                                children: "You have to be logged in to preview a template."
                            }) : (0, n.jsx)("h1", {
                                style: {color: ee},
                                children: y.error
                            }) : (0, n.jsxs)(n.Fragment, {
                                children: [(0, n.jsx)("h1", {
                                    style: {color: ee},
                                    children: "You are currently previewing a template."
                                }), (0, n.jsx)("a", {
                                    href: "/templates/".concat(y.templateId),
                                    target: "_blank",
                                    children: "Show Template"
                                })]
                            })
                        }), (0, n.jsx)("div", {
                            className: o().mainContainer, children: (0, n.jsxs)("div", {
                                className: o().userContainerWrapper,
                                style: ed,
                                children: ["simplistic" === B.layout && (0, n.jsx)("div", {
                                    className: o().profileViewsContainer,
                                    children: (0, n.jsx)(L, {
                                        publicSalt: m._PS,
                                        challenge: m._C,
                                        timestamp: m._ts,
                                        nonce: m._n,
                                        username: m.username,
                                        premiumConfig: B,
                                        views: I.page_views,
                                        containerBorderRadius: B.border_radius,
                                        isUnfold: ef === o().unfoldStart,
                                        animationStartClass: ef,
                                        isTemplatePreview: y.templatePreview
                                    })
                                }), B.parallax_animation ? (0, n.jsxs)(ec.A, {
                                    tiltEnable: !W,
                                    transitionSpeed: 700,
                                    perspective: 1e3,
                                    tiltMaxAngleX: 10,
                                    tiltMaxAngleY: 10,
                                    style: {
                                        position: "simplistic" === B.layout ? "relative" : "fixed",
                                        width: "100%",
                                        alignItems: "center",
                                        maxWidth: "var(--containerWidth)",
                                        display: "flex",
                                        justifyContent: "center",
                                        flexDirection: "column",
                                        gap: "15px"
                                    },
                                    children: [(0, n.jsx)("div", {
                                        className: ex[B.layout],
                                        style: {position: "initial"},
                                        ref: G,
                                        children: "default" === B.layout ? (0, n.jsx)(E, {
                                            userData: m,
                                            config: I,
                                            premiumConfig: B,
                                            presenceInformation: w,
                                            animationStartClass: ef,
                                            isTemplatePreview: y.templatePreview,
                                            discordLoading: S
                                        }) : "modern" === B.layout ? (0, n.jsx)(J, {
                                            userData: m,
                                            config: I,
                                            premiumConfig: B,
                                            presenceInformation: w,
                                            animationStartClass: ef,
                                            isTemplatePreview: y.templatePreview,
                                            discordLoading: S
                                        }) : "simplistic" === B.layout ? (0, n.jsx)(Q, {
                                            userData: m,
                                            config: I,
                                            premiumConfig: B,
                                            presenceInformation: w,
                                            animationStartClass: ef,
                                            setAudioData: H,
                                            audioData: M,
                                            audioPlayer: P,
                                            setAudioPlayer: R,
                                            audioContainerRef: Y,
                                            setVolume: D
                                        }) : (0, n.jsx)(em, {
                                            userData: m,
                                            config: I,
                                            premiumConfig: B,
                                            presenceInformation: w,
                                            animationStartClass: ef,
                                            isTemplatePreview: y.templatePreview
                                        })
                                    }), "enabled" === I.presence && "sleek" === B.layout && (0, n.jsx)(ep, {
                                        loading: S,
                                        presenceInformation: w,
                                        font: B.font,
                                        discordUserBadges: $,
                                        discordData: m.discord,
                                        fontSize: B.font_size,
                                        animationStartClass: ef,
                                        discordWidgetRef: X
                                    }), I.audio_player && "" !== M.url && "simplistic" !== B.layout && (0, n.jsx)(K, {
                                        animationStartClass: ef,
                                        audioPlayer: P,
                                        setAudioData: H,
                                        audioData: M,
                                        setAudioPlayer: R,
                                        audioContainerRef: Y,
                                        currentFont: B.font,
                                        audioList: I.audio,
                                        setVolume: D,
                                        fontSize: B.font_size
                                    })]
                                }) : (0, n.jsxs)("div", {
                                    className: o().usernameContainerFlex,
                                    style: {position: "simplistic" === B.layout ? "relative" : "fixed"},
                                    children: [(0, n.jsx)("div", {
                                        className: ex[B.layout],
                                        ref: G,
                                        children: "default" === B.layout ? (0, n.jsx)(E, {
                                            userData: m,
                                            config: I,
                                            premiumConfig: B,
                                            presenceInformation: w,
                                            animationStartClass: ef,
                                            isTemplatePreview: y.templatePreview,
                                            discordLoading: S
                                        }) : "modern" === B.layout ? (0, n.jsx)(J, {
                                            userData: m,
                                            config: I,
                                            premiumConfig: B,
                                            presenceInformation: w,
                                            animationStartClass: ef,
                                            isTemplatePreview: y.templatePreview,
                                            discordLoading: S
                                        }) : "simplistic" === B.layout ? (0, n.jsx)(Q, {
                                            userData: m,
                                            config: I,
                                            premiumConfig: B,
                                            presenceInformation: w,
                                            animationStartClass: ef,
                                            setAudioData: H,
                                            audioData: M,
                                            audioPlayer: P,
                                            setAudioPlayer: R,
                                            audioContainerRef: Y,
                                            setVolume: D
                                        }) : (0, n.jsx)(em, {
                                            userData: m,
                                            config: I,
                                            premiumConfig: B,
                                            presenceInformation: w,
                                            animationStartClass: ef,
                                            isTemplatePreview: y.templatePreview
                                        })
                                    }), "enabled" === I.presence && "sleek" === B.layout && (0, n.jsx)(ep, {
                                        loading: S,
                                        presenceInformation: w,
                                        font: B.font,
                                        discordUserBadges: $,
                                        discordData: m.discord,
                                        fontSize: B.font_size,
                                        animationStartClass: ef,
                                        discordWidgetRef: X
                                    }), I.audio_player && "" !== M.url && "simplistic" !== B.layout && (0, n.jsx)(K, {
                                        animationStartClass: ef,
                                        audioPlayer: P,
                                        setAudioData: H,
                                        audioData: M,
                                        setAudioPlayer: R,
                                        audioContainerRef: Y,
                                        currentFont: B.font,
                                        audioList: I.audio,
                                        setVolume: D,
                                        fontSize: B.font_size
                                    })]
                                })]
                            })
                        })]
                    })]
                }) : (null === (d = m.error) || void 0 === d ? void 0 : d.toLowerCase()) === "user not found" ? (0, n.jsx)(en, {username: k}) : m.error ? (0, n.jsx)(ea, {}) : !0 === m.ban ? (0, n.jsx)(er, {reason: m.ban_reason}) : !m.verified && (0, n.jsx)(eo, {})]
            })
        }
    }, 90367: e => {
        e.exports = {
            userInformation: "kVz23l_d0dc3a8bbeb8294cbc40",
            avatar: "zS021N_229bedc9b07046bf61d1",
            userDescription: "hSmMjc_b7ae81993acc636a36eb",
            userDescriptionText: "xc30y5_1fe480aaf30e20184c29",
            userDescriptionUsername: "rX0xXy_f8ef53b912061d92eab6",
            usernameCursor: "db1F3c_a352caa303fa59238dee",
            badgeContainerWrapper: "pwVMeN_542a0e4b8f8e3f088ce2",
            banner: "gL900l_db8fee3064281101ea2a",
            adjustBanner: "zbzENV_5531de12e00b53c729b8",
            socialsWrapper: "d20bYY_b4c509d547ebc284a4d0"
        }
    }, 97419: (e, t, s) => {
        "use strict";
        s.d(t, {A: () => r, k: () => a});
        var n = s(12115);
        let a = {
            satoshi: "'Satoshi', sans-serif",
            firacode: "'Fira Code', monospace",
            poppins: "'Poppins', sans-serif",
            array: "'Array', sans-serif",
            drippy: "'Were-Beast', sans-serif",
            minecraft: "'Minecraftia', sans-serif",
            outfit: "'Outfit', sans-serif",
            chillax: "'Chillax', sans-serif",
            velocity: "'Panchang', sans-serif",
            technor: "'Technor', sans-serif",
            jetbrains_mono: "'JetBrains Mono', monospace",
            silkscreen: "'Silkscreen', serif",
            parkinsans: "'Parkinsans', serif"
        };

        function r(e) {
            let {font: t, letterSpacing: s} = e, r = {1: "-0.5px", 2: null, 3: "1.5px"};
            return (0, n.useEffect)(() => {
                let e = document.head;
                if ("object" == typeof t) {
                    let e = t.url;
                    new FontFace("USERFONT", "url(".concat(e, ")")).load().then(e => {
                        document.fonts.add(e), document.body.style.fontFamily = "USERFONT, sans-serif";
                        let t = s ? r[s] : null;
                        t && (document.body.style.letterSpacing = t)
                    }).catch(e => {
                        console.error("Failed to load font:", e)
                    })
                } else {
                    let s = document.createElement("link");
                    switch (s.rel = "stylesheet", s.dataset.font = t, t) {
                        case"firacode":
                            s.href = "https://fonts.googleapis.com/css2?family=Fira+Code:wght@500&display=swap";
                            break;
                        case"poppins":
                            s.href = "https://fonts.googleapis.com/css2?family=Poppins&display=swap";
                            break;
                        case"array":
                            s.href = "https://api.fontshare.com/v2/css?f[]=array@400&display=swap";
                            break;
                        case"drippy":
                            s.href = "https://fonts.cdnfonts.com/css/were-beast-2";
                            break;
                        case"minecraft":
                            s.href = "https://fonts.cdnfonts.com/css/minecraftia";
                            break;
                        case"outfit":
                            s.href = "https://fonts.googleapis.com/css2?family=Outfit:wght@100..900&display=swap";
                            break;
                        case"chillax":
                            s.href = "https://api.fontshare.com/v2/css?f[]=chillax@500&display=swap";
                            break;
                        case"velocity":
                            s.href = "https://api.fontshare.com/v2/css?f[]=panchang@500&display=swap";
                            break;
                        case"technor":
                            s.href = "https://api.fontshare.com/v2/css?f[]=technor@400&display=swap";
                            break;
                        case"jetbrains_mono":
                            s.href = "https://fonts.googleapis.com/css2?family=JetBrains+Mono:ital@0;1&display=swap";
                            break;
                        case"silkscreen":
                            s.href = "https://fonts.googleapis.com/css2?family=Silkscreen&display=swap";
                            break;
                        case"parkinsans":
                            s.href = "https://fonts.googleapis.com/css2?family=Parkinsans:wght@500&display=swap";
                            break;
                        default:
                            s.href = ""
                    }
                    s.href && e.appendChild(s);
                    let n = document.createElement("style");
                    return n.dataset.font = t, n.innerHTML = "\n                body, * {\n                    font-family: ".concat(a[t] || "sans-serif", " !important;\n                }\n            "), e.appendChild(n), () => {
                        s.parentNode && e.removeChild(s), n.parentNode && e.removeChild(n)
                    }
                }
            }, [t]), null
        }
    }
}]);