/**
 * Copyright (c) 2020-2021 Weitian Leung
 *
 * This file is part of pywpsrpc.
 *
 * This file is distributed under the MIT License.
 * See the LICENSE file for details.
 *
*/

namespace etapi
{
    struct IFormatConditions : public IDispatch /Abstract/
    {
    public:
        virtual HRESULT get_Application(
            Application **RHS /Out/) = 0;

        virtual HRESULT get_Creator(
            XlCreator *RHS /Out/) = 0;

        virtual HRESULT get_Count(
            long *RHS /Out/) = 0;

        Py_ssize_t __len__() const;
        %MethodCode
            long count = 0;
            if (sipCpp->get_Count(&count) != S_OK)
                count = 0;
            sipRes = count;
        %End

        virtual HRESULT Item(
            VARIANT Index,
            IFormatCondition **RHS /Out/) = 0;
        %MethodCode
            sipRes = sipCpp->Item(*a0, (IDispatch**)&a1);
        %End

        IFormatCondition* __getitem__(VARIANT index) const;
        %MethodCode
            etapi::IFormatCondition *prop = nullptr;
            if (sipCpp->Item(*a0, (IDispatch**)&prop) != S_OK)
                sipIsErr = 1;
            else
                sipRes = prop;
        %End

        virtual HRESULT Add(
            XlFormatConditionType Type,
            VARIANT Operator = argMissing2(),
            VARIANT Formula1 = argMissing2(),
            VARIANT Formula2 = argMissing2(),
            VARIANT String = argMissing2(),
            VARIANT TextOperator = argMissing2(),
            VARIANT DateOperator = argMissing2(),
            VARIANT ScopeType = argMissing2(),
            IFormatCondition **RHS /Out/ = 0) = 0;
        %MethodCode
            sipRes = sipCpp->Add(a0, *a1, *a2, *a3, *a4, *a5, *a6, *a7, (IDispatch**)&a8);
        %End

        virtual HRESULT get__Default(
            VARIANT Index,
            IDispatch **RHS /Out/) = 0;

        virtual HRESULT Delete() = 0;

        virtual HRESULT AddColorScale(
            long ColorScaleType,
            IDispatch **RHS /Out/) = 0;

        virtual HRESULT AddDatabar(
            IDispatch **RHS /Out/) = 0;

        virtual HRESULT AddIconSetCondition(
            IDispatch **RHS /Out/) = 0;

        virtual HRESULT AddTop10(
            IDispatch **RHS /Out/) = 0;

        virtual HRESULT AddAboveAverage(
            AboveAverage **RHS /Out/) = 0;
        %MethodCode
            sipRes = sipCpp->AddAboveAverage((IDispatch**)&a0);
        %End

        virtual HRESULT AddUniqueValues(
            IDispatch **RHS /Out/) = 0;

        Application* getApplication();
        %MethodCode
            etapi::Application *prop = nullptr;
            HRESULT hr = sipCpp->get_Application(&prop);
            if (hr != S_OK)
            {
                PyErr_Format(PyExc_AttributeError,
                    "Call 'get_Application()' failed with 0x%x", hr);
                return nullptr;
            }
            else
            {
                sipRes = prop;
            }
        %End

        %Property(name=Application, get=getApplication)

        XlCreator getCreator();
        %MethodCode
            etapi::XlCreator prop = (etapi::XlCreator)0;
            HRESULT hr = sipCpp->get_Creator(&prop);
            if (hr != S_OK)
            {
                PyErr_Format(PyExc_AttributeError,
                    "Call 'get_Creator()' failed with 0x%x", hr);
                return nullptr;
            }
            else
            {
                sipRes = prop;
            }
        %End

        %Property(name=Creator, get=getCreator)

        long getCount();
        %MethodCode
            long prop = 0;
            HRESULT hr = sipCpp->get_Count(&prop);
            if (hr != S_OK)
            {
                PyErr_Format(PyExc_AttributeError,
                    "Call 'get_Count()' failed with 0x%x", hr);
                return nullptr;
            }
            else
            {
                sipRes = prop;
            }
        %End

        %Property(name=Count, get=getCount)

    };
};
