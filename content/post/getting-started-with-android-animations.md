+++
description = "A quick guide to the basics of Android animations"
title = "Getting Started with Android Animations"
date = "2017-07-03T09:19:07-06:00"
draft = true
categories = ["android"]
tags = ["android", "animations"]

+++



Android Animations has come a long way from Froyo to Honeycomb and then to Marshmallow. This has included a lot of API changes in Android. The truth - there is no animation library that you can just include in your app and they start automagically working. Animations need to be coded. There is no way around it. Even though, animateLayoutChanges="true" has some nice usages, it's likely rare.

When Material Design became the standard, Android went all out on colors and motion, moreover, delighting the user. Animations in Android have a complex vocabulary to learn, even though it has gotten easier over the years, it is still hefty. I might not even cover all the vocabulary in this post however it will give you a core idea of how to get started with animations on Android to delight your users. Android Docs are great! No doubt, however they are usually jargon to start out with.

## View Animations

#### Animating a button

Views are animated with the ViewPropertyAnimator. With the fluent API you can just do view.animate() to start animating your views.

```
btnAnimate.animate().translationY(500).alpha(0).setDuration(1000)
                .setListener(new AnimatorListenerAdapter() {
            @Override
            public void onAnimationCancel(Animator animation) {
                super.onAnimationCancel(animation);
                btnAnimate.setVisibility(View.GONE);
            }
        });
// They run in parallel, alpha and translationY

```
Don't override AnimatorListener, override AnimatorListenerAdapter. Try it!

I have seen developers at Google recommend this method as to the following one as it makes animations fluent.

#### ObjectAnimator

The other method is to use the ObjectAnimator.

ObjectAnimator fade = ObjectAnimator.ofFloat(btnAnimate, View.ALPHA, 0);
fade.start();
Interpolators

#### Linear vs Curve

Animations need to look real to your user. A decelerating/accelerating object always looks better than a object in linear motion.

Interpolator is another class in the API which makes it easy to define the rate of change of animations. These allow view properties like alpha and translate to be accelerated and decelerated (it does a lot more).

```java
  btnAnimate.animate()
                .setInterpolator(new AccelerateDecelerateInterpolator())
                .setDuration(2500)
                .alpha(0).translationY(500);
```
The interpolator changes the alpha of the view by accelerating the fade in the beginning and decelerating completely disappears from the screen. The same with translationY.

#### Sets of Animations

Translate and Fade Off Screen

Your might be thinking, I can just do this, however animations run asynchronously hence the following code will run in parallel.

```java

// First Translate
        btnAnimate.animate().translationY(400);
// Then Fade Out
        btnAnimate.animate().alpha(0);
// They both should run sequentially

```
In comes the AnimationSet. This allows animations to be played sequentially, in parallel and a lot more. Basically, AnimationSet is the choreographer of your animations.

With AnimationSet you can no longer use the view.animate() but now are using the ObjectAnimator.
```

        ObjectAnimator fade = ObjectAnimator.ofFloat(btnAnimate, View.ALPHA, 0);
        ObjectAnimator translate = ObjectAnimator.ofFloat(btnAnimate, View.TRANSLATION_Y, 400);
        AnimatorSet set = new AnimatorSet();
        set.playSequentially(translate, fade);
        set.start();
      // The button will translate and then fade out.
      // This example is just scratching the surface.
```

In the end...

With ViewPropertyAnimator, Interpolators and AnimatorSet you can choreograph many animations that delight the users with the use of motion. ViewPropertyAnimator has a lot of advantages of its predecessor ObjectAnimator, namely the fluent API. However ObjectAnimator still is the mother of all while choreographing complex animations. Also, don't forget interpolators. Its these that make your animations delight the user.

```
NOTE: This is a article I had on my old blog that I ported over to this version.
```
